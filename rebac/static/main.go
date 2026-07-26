package main

import (
	"fmt"
	"strings"
)

type Tuple struct {
	Object   string
	Relation string
	User     string
}

type UsersetRef struct {
	Object   string
	Relation string
}

type TupleStore interface {
	Read(object, relation string) ([]Tuple, error)
}

type RewriteRule struct {
	This            bool
	ComputedUserSet string
	TupleToUserSet  *T2U
	Union           []RewriteRule
	Intersection    []RewriteRule
	Exclusion       *ExclusionRule
}

type T2U struct {
	TuplesetRelation string
	ComputedRelation string
}

type ExclusionRule struct {
	Base    RewriteRule
	Exclude RewriteRule
}

type NameSpace struct {
	Name     string
	Relation map[string]RewriteRule
}

func DocNamespace() NameSpace {
	return NameSpace{
		Name: "doc",
		Relation: map[string]RewriteRule{
			"owner": {This: true},
			"editor": {
				Union: []RewriteRule{
					{This: true},
					{ComputedUserSet: "owner"},
				},
			},
			"viewer": {
				Union: []RewriteRule{{This: true},
					{ComputedUserSet: "editor"},
					{TupleToUserSet: &T2U{
						TuplesetRelation: "parent",
						ComputedRelation: "viewer",
					}},
				},
			},
		},
	}
}

type Checker struct {
	Store     TupleStore
	Namespace NameSpace
}

func (c *Checker) Check(object, relation, user string) (bool, error) {
	fmt.Printf("Relation: %+v\n", c.Namespace.Relation)
	rule, ok := c.Namespace.Relation[relation]
	fmt.Printf("This is rule for relation:%s: %+v\n", relation, rule)
	if !ok {
		return false, fmt.Errorf("Unknown relation %q", relation)
	}
	return c.evalRule(object, relation, rule, user)
}

func (c *Checker) evalRule(object, relation string, rule RewriteRule, user string) (bool, error) {
	defer fmt.Println("Returning from eval rule")
	switch {
	case rule.This:
		fmt.Printf("Yes rule.this\n")
		val, err := c.checkDirect(object, relation, user)
		fmt.Printf("value: %t, err: %v\n", val, err)
		return val, err
	case rule.ComputedUserSet != "":
		return c.Check(object, rule.ComputedUserSet, user)

	case rule.TupleToUserSet != nil:
		return c.evalTupleToUserset(object, *rule.TupleToUserSet, user)

	case rule.Union != nil:
		fmt.Printf("Yes rule.Union\n")
		for _, child := range rule.Union {
			ok, err := c.evalRule(object, relation, child, user)
			if err != nil {
				return false, err
			}
			if ok {
				return true, nil
			}
		}
		return false, nil
	case rule.Intersection != nil:
		for _, rule := range rule.Intersection {
			ok, err := c.evalRule(object, relation, rule, user)
			if err != nil || !ok {
				return false, err
			}
		}
		return true, nil
	case rule.Exclusion != nil:
		base, err := c.evalRule(object, relation, rule.Exclusion.Base, user)
		if err != nil || !base {
			return false, err
		}
		excluded, err := c.evalRule(object, relation, rule.Exclusion.Exclude, user)
		if err != nil {
			return false, err
		}
		return !excluded, nil
	}
	return false, nil
}

func (c *Checker) evalTupleToUserset(object string, t2u T2U, user string) (bool, error) {
	tuples, err := c.Store.Read(object, t2u.TuplesetRelation)
	if err != nil {
		return false, err
	}
	for _, t := range tuples {
		ok, err := c.Check(t.User, t2u.ComputedRelation, user)
		if err != nil {
			return false, err
		}
		if ok {
			return true, nil
		}
	}
	return false, nil
}
func (c *Checker) checkDirect(object, relation, user string) (bool, error) {
	tuples, err := c.Store.Read(object, relation)
	fmt.Printf("This is tuples: %s\n", tuples)
	if err != nil {
		return false, err
	}
	for _, t := range tuples {
		fmt.Println("Inside loo in checkDir")
		if t.User == user {
			return true, nil
		}
		if ref, isUserset := parseUsersetRef(t.User); isUserset {
			ok, err := c.Check(ref.Object, ref.Relation, user)
			if err != nil || ok {
				return ok, err
			}
		}
	}
	return false, nil
}

func parseUsersetRef(user string) (UsersetRef, bool) {
	idx := strings.Index(user, "#")
	if idx == -1 {
		return UsersetRef{}, false
	}
	return UsersetRef{
		Object:   user[:idx],   //folder:eng
		Relation: user[idx+1:], //viewer
	}, true
}

type InMemoryStore struct {
	tuples []Tuple
}

func (s *InMemoryStore) Read(object, relation string) ([]Tuple, error) {
	var result []Tuple
	for _, t := range s.tuples {
		fmt.Printf("Tuples: %+v\n", t)
		if t.Object == object && t.Relation == relation {
			result = append(result, t)
		}
	}
	return result, nil
}

func main() {

	store := &InMemoryStore{tuples: []Tuple{
		{Object: "doc:readme", Relation: "owner", User: "nishan"},
		{Object: "doc:readme", Relation: "parent", User: "folder:eng"},
		{Object: "folder:eng", Relation: "viewer", User: "hehe"},
	}}

	checker := &Checker{Store: store, Namespace: DocNamespace()}

	ok, _ := checker.Check("doc:readme", "owner", "nishan")
	fmt.Println(ok)
}
