Counter(
    in Clk,
    out reg [8] Counter@Clk)
{
    Counter <= Counter + 1
}
