module Adder(
    input [7:0] A,
    input [7:0] B,
    output [7:0] C
);
    assign C = A + B;
endmodule

module Shifter(
    input [7:0] A,
    input [7:0] B,
    output [7:0] C
);
    assign C = A << B;
endmodule

module AddSub(
    input [7:0] Val,
    output [7:0] Add,
    output [7:0] Sub
);
    assign Add = Val + 1;
    assign Sub = Val - 1;
endmodule

