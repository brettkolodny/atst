#!/usr/bin/env -S deno run --allow-write

Deno.addSignalListener("SIGINT", () => {
  console.log("Goodbye!");
  Deno.exit(1);
});

// const countDown = async () => {
//   let count = parseInt(Deno.args[0] ?? "10");

//   while (count > 0) {
//     console.log(`Blast off in ${count}...`);
//     await new Promise((r) => setTimeout(r, 1000));
//     count -= 1;
//   }

//   console.log("Blast off! 🚀");
// };

// const countUp = async () => {
//   // TODO
// };

const main = async () => {
  let count = parseInt(Deno.args[0] ?? "10");

  while (count > 0) {
    console.log(`Blast off in ${count}...`);
    await new Promise((r) => setTimeout(r, 1000));
    count -= 1;
  }

  console.log("Blast off! 🚀");
  Deno.exit(1);
};

await main();
