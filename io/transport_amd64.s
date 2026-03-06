//go:build linux && amd64
// +build linux,amd64

#include "textflag.h"

// Heimdall_Recv(fd int, msgbatch uintptr, vlen int) (n int, err int)
// ABIInternal: AX=fd, BX=msgbatch, CX=vlen
TEXT ·Heimdall_Recv(SB), NOSPLIT, $0-40
    MOVQ	AX, DI
    MOVQ	BX, SI
    MOVQ	CX, DX
    MOVQ	$0, R10
    MOVQ	$0, R8
retry:
    MOVQ	$299, AX // SYS_RECVMMSG
    SYSCALL
    CMPQ	AX, $-4    // EINTR
    JE	    retry
    
    CMPQ	AX, $-4095
    JAE	    err_label
    
    // Success: n is in AX, set err (BX) to 0
    XORQ	BX, BX     // err = 0
    MOVQ	AX, ret+24(FP)
    MOVQ	BX, ret1+32(FP)
    RET

err_label:
    NEGQ	AX
    MOVQ    AX, BX    // err = errno
    MOVQ    $-1, AX   // n = -1
    MOVQ	AX, ret+24(FP)
    MOVQ	BX, ret1+32(FP)
    RET

// Heimdall_Send(fd int, msgbatch uintptr, vlen int) (n int, err int)
TEXT ·Heimdall_Send(SB), NOSPLIT, $0-40
    MOVQ	AX, DI
    MOVQ	BX, SI
    MOVQ	CX, DX
    MOVQ	$0, R10
retry_s:
    MOVQ	$307, AX // SYS_SENDMMSG
    SYSCALL
    CMPQ	AX, $-4
    JE	    retry_s
    
    // Check error
    CMPQ	AX, $-4095
    JAE	    err_send_label
    
    // Success: n is in AX, set err (BX) to 0
    XORQ	BX, BX     // err = 0
    MOVQ	AX, ret+24(FP)
    MOVQ	BX, ret1+32(FP)
    RET

err_send_label:
    NEGQ	AX
    MOVQ    AX, BX    // err = errno
    MOVQ    $-1, AX   // n = -1
    MOVQ	AX, ret+24(FP)
    MOVQ	BX, ret1+32(FP)
    RET
