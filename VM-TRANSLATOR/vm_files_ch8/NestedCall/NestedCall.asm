// Boostrap code
@256
D=A
@SP
M=D
// call Sys.init 0
@return_address_0
D=A
@SP
A=M
M=D
@SP
M=M+1
@LCL
D=M
@SP
A=M
M=D
@SP
M=M+1
@ARG
D=M
@SP
A=M
M=D
@SP
M=M+1
@THIS
D=M
@SP
A=M
M=D
@SP
M=M+1
@THAT
D=M
@SP
A=M
M=D
@SP
M=M+1

D=M
@5
D=D-A
@ARG
M=D

@SP
D=M
@LCL
M=D
@Sys.init
0;JMP
(return_address_0)
0;JMP
// function Sys.init0
(Sys.init)
@SP
A=M
D=A
@SP
M=D
// C_PUSH constant[0]
@0
D=A
@SP
A=M
M=D
@SP
M=M+1
// C_POP pointer[0]
@SP
AM=M-1
D=M
@THIS
M=D
// C_PUSH constant[5000]
@5000
D=A
@SP
A=M
M=D
@SP
M=M+1
// C_POP pointer[1]
@SP
AM=M-1
D=M
@THAT
M=D
// call Sys.main 0
@return_address_1
D=A
@SP
A=M
M=D
@SP
M=M+1
@LCL
D=M
@SP
A=M
M=D
@SP
M=M+1
@ARG
D=M
@SP
A=M
M=D
@SP
M=M+1
@THIS
D=M
@SP
A=M
M=D
@SP
M=M+1
@THAT
D=M
@SP
A=M
M=D
@SP
M=M+1

D=M
@5
D=D-A
@ARG
M=D

@SP
D=M
@LCL
M=D
@Sys.main
0;JMP
(return_address_1)
// C_POP temp[1]
@1
D=A
@R5
D=D+A
@FRAME
M=D
@SP
AM=M-1
D=M
@FRAME
A=M
M=D
// label LOOP
(LOOP)
// goto LOOP
@LOOP
0;JMP
// function Sys.main5
(Sys.main)
@SP
A=M
M=0
A=A+1
M=0
A=A+1
M=0
A=A+1
M=0
A=A+1
M=0
A=A+1
D=A
@SP
M=D
// C_PUSH constant[4001]
@4001
D=A
@SP
A=M
M=D
@SP
M=M+1
// C_POP pointer[0]
@SP
AM=M-1
D=M
@THIS
M=D
// C_PUSH constant[5001]
@5001
D=A
@SP
A=M
M=D
@SP
M=M+1
// C_POP pointer[1]
@SP
AM=M-1
D=M
@THAT
M=D
// C_PUSH constant[200]
@200
D=A
@SP
A=M
M=D
@SP
M=M+1
// C_POP local[1]
@1
D=A
@LCL
D=D+M
@FRAME
M=D
@SP
AM=M-1
D=M
@FRAME
A=M
M=D
// C_PUSH constant[40]
@40
D=A
@SP
A=M
M=D
@SP
M=M+1
// C_POP local[2]
@2
D=A
@LCL
D=D+M
@FRAME
M=D
@SP
AM=M-1
D=M
@FRAME
A=M
M=D
// C_PUSH constant[6]
@6
D=A
@SP
A=M
M=D
@SP
M=M+1
// C_POP local[3]
@3
D=A
@LCL
D=D+M
@FRAME
M=D
@SP
AM=M-1
D=M
@FRAME
A=M
M=D
// C_PUSH constant[123]
@123
D=A
@SP
A=M
M=D
@SP
M=M+1
// call Sys.add12 1
@return_address_2
D=A
@SP
A=M
M=D
@SP
M=M+1
@LCL
D=M
@SP
A=M
M=D
@SP
M=M+1
@ARG
D=M
@SP
A=M
M=D
@SP
M=M+1
@THIS
D=M
@SP
A=M
M=D
@SP
M=M+1
@THAT
D=M
@SP
A=M
M=D
@SP
M=M+1

D=M
@6
D=D-A
@ARG
M=D

@SP
D=M
@LCL
M=D
@Sys.add12
0;JMP
(return_address_2)
// C_POP temp[0]
@0
D=A
@R5
D=D+A
@FRAME
M=D
@SP
AM=M-1
D=M
@FRAME
A=M
M=D
// C_PUSH local[0]
@0
D=A
@LCL
A=D+M
D=M
@SP
A=M
M=D
@SP
M=M+1
// C_PUSH local[1]
@1
D=A
@LCL
A=D+M
D=M
@SP
A=M
M=D
@SP
M=M+1
// C_PUSH local[2]
@2
D=A
@LCL
A=D+M
D=M
@SP
A=M
M=D
@SP
M=M+1
// C_PUSH local[3]
@3
D=A
@LCL
A=D+M
D=M
@SP
A=M
M=D
@SP
M=M+1
// C_PUSH local[4]
@4
D=A
@LCL
A=D+M
D=M
@SP
A=M
M=D
@SP
M=M+1
// add
@SP
AM=M-1
D=M
@SP
AM=M-1
M=D+M
@SP
M=M+1
// add
@SP
AM=M-1
D=M
@SP
AM=M-1
M=D+M
@SP
M=M+1
// add
@SP
AM=M-1
D=M
@SP
AM=M-1
M=D+M
@SP
M=M+1
// add
@SP
AM=M-1
D=M
@SP
AM=M-1
M=D+M
@SP
M=M+1
// return
@LCL
D=M
@FRAME
M=D
@FRAME
D=M
@5
D=D-A
A=D
D=M
@RET
M=D
@SP
AM=M-1
D=M
@ARG
A=M
M=D
@ARG
D=M+1
@SP
M=D
@FRAME
D=M
@1
D=D-A
A=D
D=M
@THAT
M=D
@FRAME
D=M
@2
D=D-A
A=D
D=M
@THIS
M=D
@FRAME
D=M
@3
D=D-A
A=D
D=M
@ARG
M=D
@FRAME
D=M
@4
D=D-A
A=D
D=M
@LCL
M=D
@RET
A=M
0;JMP
// function Sys.add120
(Sys.add12)
@SP
A=M
D=A
@SP
M=D
// C_PUSH constant[4002]
@4002
D=A
@SP
A=M
M=D
@SP
M=M+1
// C_POP pointer[0]
@SP
AM=M-1
D=M
@THIS
M=D
// C_PUSH constant[5002]
@5002
D=A
@SP
A=M
M=D
@SP
M=M+1
// C_POP pointer[1]
@SP
AM=M-1
D=M
@THAT
M=D
// C_PUSH argument[0]
@0
D=A
@ARG
A=D+M
D=M
@SP
A=M
M=D
@SP
M=M+1
// C_PUSH constant[12]
@12
D=A
@SP
A=M
M=D
@SP
M=M+1
// add
@SP
AM=M-1
D=M
@SP
AM=M-1
M=D+M
@SP
M=M+1
// return
@LCL
D=M
@FRAME
M=D
@FRAME
D=M
@5
D=D-A
A=D
D=M
@RET
M=D
@SP
AM=M-1
D=M
@ARG
A=M
M=D
@ARG
D=M+1
@SP
M=D
@FRAME
D=M
@1
D=D-A
A=D
D=M
@THAT
M=D
@FRAME
D=M
@2
D=D-A
A=D
D=M
@THIS
M=D
@FRAME
D=M
@3
D=D-A
A=D
D=M
@ARG
M=D
@FRAME
D=M
@4
D=D-A
A=D
D=M
@LCL
M=D
@RET
A=M
0;JMP