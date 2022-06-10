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
The flagship feature of this language.

One of the hardest parts of FPGA development is developing state machines and getting the timing of signals correct. Chrono aims to assist in designing these by introducing a feature called supporting sequences.

This is a feature is inspired by SystemVerilog's ability to sequence things in an initial/always block using delays, wait events, and process forking. While extremely powerful, this is not synthesize and is just limited to test benches.

Sequences allows you program a state machine or sequence of signals as intuitively as possible. Instead of having to manually create the state machine, Chrono will automatically generate a state machine that will produce the sequence of events that you describe.

For this feature, Chrono uses a sequential/parallel paradigm. By default all statements are parallel and operate the way you're used to in HDL languages. However, you're also able to declare a sequential block where each statement operates is executed one at a time after every clock.

The powerful comes from the fact that these parallel and sequential blocks can be nested inside each other. Therefore you can have a sequence of nested sequences running in parallel!

This should make coding complex operations and interfaces *significantly* easier than it is to do manually in VHDL and Verilog.

</br>

## Design Philosophy

### No fluff
Chrono should be as concise as possible. Requiring enough explicit syntax to ensure you know you're telling the tools to do, without being unnecessarily wordy.

Therefore we use the syntax `[8]` to create a signal of 8 bits wide, instead of `(7 downto 0)` in the case of VHDL or `[7:0]` in the case of verilog.

</br>

### The right way should be the easiest way
The language should naturally assist you in creating circuits correctly. This makes it less likely that you'll run into issues during design.

In FPGAs, resets are normally not necessary. Registers are initialized into a known state by the GSR (Global Set/Reset). This means adding large reset nets to your circuit is actually creating a second redundant global reset. Large resets can take up priority routing within the FPGA to ensure it meets timing on the high fanout net. This should be avoided when unnecessary. Therefore specifying default initialization values should be very easy while adding an explicit reset should require extra code.








