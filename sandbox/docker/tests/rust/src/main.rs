pub fn hello() -> &'static str {
    "Hello, World!"
}

fn main() {
    let hello_world = hello();
    print!("{}", hello_world);
}
