Mux(
    in [8] A,
    in [8] B,
    in [8] C,
    out [8] Q,
    in [2] Sel)
{
    if (Sel == 0)
    {
        Q = A
    }
    else if (Sel == 1)
    {
        Q = B
    }
    else if (Sel == 2)
    {
        Q = C
    }
    else
    {
        Q = 0
    }
}
