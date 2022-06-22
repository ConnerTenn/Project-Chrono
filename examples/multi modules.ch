
Adder(in [8] A, in [8] B, out [8] C)
{
    C = A + B - C;
}

Shifter(in [8] A, in [8] B, out [8] C)
{
    C = A << B * 2;
}

AddSub(in [8] Val, out [8] Add, out [8] Sub)
{
    Add = Val + 1 * 2;
    Sub = (Val - 1) * 2;
}

