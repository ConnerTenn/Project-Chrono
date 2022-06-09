module Adder(
    In [7:0] A,
    In [7:0] B,
    Out [7:0] C
);
    assign C = A + B;
endmodule

module Shifter(
    In [7:0] A,
    In [7:0] B,
    Out [7:0] C
);
    assign C = A << B;
endmodule

module AddSub(
    In [7:0] Val,
    Out [7:0] Add,
    Out [7:0] Sub
);
    assign Add = Val + 1;
    assign Sub = Val - 1;
endmodule

