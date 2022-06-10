# Project-Chrono

[Check out the wiki for documentation](https://github.com/ConnerTenn/Project-Chrono/wiki)

Note: This is still in very early development and much is subject to change.


## Purpose

Chrono is a new HDL language for FPGA development.

We already have VHDL and Verilog. Why add *yet another* language?

Chrono aims to be a modern take on HDL languages. VHDL and Verilog are older languages that have some quirks in their design. Chrono aims to clean up the language and make it easier develop code that does what you want, instead of having to work hard to get the tools to interpret what you want correctly.

Chrono aims to be easy to read and use, while also giving you confidence in the circuit that will be generated.

</br>

## Sequences
The flagship feature of this language!

One of the hardest parts of FPGA development is developing state machines and getting the timing of signals correct. Chrono aims to assist in designing these by introducing a feature called supporting sequences.

This is a feature is inspired by SystemVerilog's ability to sequence things in an initial/always block using delays, wait events, and process forking. While extremely powerful, this is not synthesize and is just limited to test benches.

Sequences allows you program a state machine or sequence of signals as intuitively as possible. Instead of having to manually create the state machine, Chrono will automatically generate a state machine that will produce the sequence of events that you describe.

For this feature, Chrono uses a sequential/parallel paradigm. By default all statements are parallel and operate the way you're used to in HDL languages. However, you're also able to declare a sequential block where each statement operates is executed one at a time after every clock.

The powerful comes from the fact that these parallel and sequential blocks can be nested inside each other. Therefore you can have a sequence of nested sequences running in parallel!

```verilog
//Sequential
@(Clk)
{
    A <= 1;
    A <= 2;

    //parallel
    {
        A <= 3;
        B <= 5;
    }

    //parallel
    {
        //Sequential
        @(clk)
        {
            A <= 4;
            A <= 5;
        }
        B <= B + 1; //Count while the sequence is running
    } //Progresses after all contained sequences are done

    A <= 0;
    B <= 0;
}

//           __    __    __    __    __    __    __  
//Clk ______/  \__/  \__/  \__/  \__/  \__/  \__/  \_
//A   0     1     2     3     4     5     0          
//B   0                 5     6     7           0    
```
[Checkout the wiki for more examples!](https://github.com/ConnerTenn/Project-Chrono/wiki/Sequences)

This should make coding complex operations and interfaces *significantly* easier than it is to do manually in VHDL and Verilog.

</br>

## Design Philosophy

### No fluff
Chrono should be as concise as possible. Requiring enough explicit syntax to ensure you know you're telling the tools to do, without being unnecessarily wordy.

Therefore curly braces `{}` are used instead of the `begin` and `end` found in other languages.

The syntax with square brackets is used to create a vector signal. For example `[8]`, instead of `(7 downto 0)` in the case of VHDL or `[7:0]` in the case of verilog.

</br>

### The right way should be the easiest way
The language should naturally assist you in creating circuits correctly. This makes it less likely that you'll run into issues during design.

In FPGAs, resets are normally not necessary. Registers are initialized into a known state by the GSR (Global Set/Reset). This means adding large reset nets to your circuit is actually creating a second redundant global reset. Large resets can take up priority routing within the FPGA to ensure it meets timing on the high fanout net. This should be avoided when unnecessary. Therefore specifying default initialization values should be very easy while adding an explicit reset should require extra code.

Latches should also be avoided when developing for FPGAs, while FPGAs do have tha ability to create latches, they are much more difficult to create timing constraints for and are prone to causing issues with setup and hold timing. In Verilog and especially VHDL, it is extremely easy to accidentally create latches. This is not desired for a language. It's still possible to create latches in Chrono, but the syntax for creating one is purposely less straight forward, to ensure it is never accidentally done.

To support this, every register is declared with a clock it is synchronous to using the `@` operator. Then synchronous assignments can be made easily in the main block using the `<=` operator.
```verilog
sig [8] Counter@(Clk)

Counter <= Counter + 1;
```

This is in contrast to VHDL which infers registers based on a processes's sensitivity list.
```vhdl
signal Counter : std_logic_vector(7 downto 0);

process(Clk)
begin
    if rising_edge(Clk) then
    begin
        Counter <= Counter + 1;
    end if;
end
```

Or Verilog, which is more explicit with the `reg` keyword, however that still has some counter intuitive behavior. [(See the next section for more detail)](#No-unintuitive-behavior-or-ambiguities)

```verilog
reg [7:0] Counter;

always @(posedge Clk)
begin
    Counter <= Counter + 1;
end
```

</br>

### No unintuitive behavior or ambiguities

The circuits generated by Chrono should be representative of what the designing intended. There shouldn't be any strange counterintuitive quirks.

For example, Verilog has the `reg` keyword, however there are instances where a signal *must* be declared as a register, however gets inferred to be a wire.

```verilog
reg [7:0] myval;

always @(option)
begin
    case (option)
        'h0: myval = 1;
        'h1: myval = 3;
        'h2: myval = 5;
        default: myval = 0;
    endcase
end
```

In this instance, myval is declared as a `reg` but actually gets inferred to be a wire with a combinatorial mux.

Issues like these should never exist in Chrono. It should always be clear what you are describing in a language. No surprises!

</br>

### Essential features should be builtin
Essential and commonly used features of the language should be built into the language. For example, finding the number of bits required to store a particular value.

