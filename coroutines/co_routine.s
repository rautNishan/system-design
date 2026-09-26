#include "textflag.h"

TEXT ·co_yeild(SB), NOSPLIT, $0-0
    MOVQ SP, DI
    JMP ·test(SB)
    
