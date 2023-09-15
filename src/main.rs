mod python;
mod ruby;

fn main() {
    ruby::dump_ruby_structs_ruby_2_6_0();
    ruby::dump_ruby_structs_ruby_2_6_3();

    ruby::dump_ruby_structs_ruby_2_7_1();
    ruby::dump_ruby_structs_ruby_2_7_4();
    ruby::dump_ruby_structs_ruby_2_7_6();

    ruby::dump_ruby_structs_ruby_3_0_0();
    ruby::dump_ruby_structs_ruby_3_0_4();
    ruby::dump_ruby_structs_ruby_3_1_2();
    ruby::dump_ruby_structs_ruby_3_1_3();
    ruby::dump_ruby_structs_ruby_3_2_0();
    ruby::dump_ruby_structs_ruby_3_2_1();

    python::dump_python_structs_2_7_15();

    python::dump_python_structs_3_3_7();
    python::dump_python_structs_3_4_8();
    python::dump_python_structs_3_5_5();
    python::dump_python_structs_3_6_6();
    python::dump_python_structs_3_7_0();

    python::dump_python_structs_3_8_0();
    python::dump_python_structs_3_9_5();
    python::dump_python_structs_3_10_0();
    python::dump_python_structs_3_11_0();
}
