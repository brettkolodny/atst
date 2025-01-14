#!/usr/bin/env -S deno run --allow-write

Deno.addSignalListener("SIGINT", () => {
  console.log("Launch aborted!");
  Deno.exit(1);
});

const countDown = async () => {
  let count = 10;

  while (count > 0) {
    console.log(`Blast off in ${count}...`);
    await new Promise((r) => setTimeout(r, 1000));
    count -= 1;
  }

  console.log("Blast off! 🚀");
  Deno.exit(1);
};

const countUp = async () => {
  let count = 0;

  while (++count) {
    console.log(`The count is now ${count}...`);
    await new Promise((r) => setTimeout(r, 1000));
  }
};

const main = async () => {
  if (Deno.args[0] === "countdown") {
    await countDown();
  } else {
    await countUp();
  }
};

await main();
