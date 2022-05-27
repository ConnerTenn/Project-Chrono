
Adder(in [8] A, in [8] B, out [8] C)
{
    C = A + B;
}

Shifter(in [8] A, in [8] B, out [8] C)
{
    C = A << B;
}

AddSub(in [8] Val, out [8] Add, out [8] Sub)
{
    Add = Val + 1;
    Sub = Val - 1;
}

