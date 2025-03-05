#!/usr/bin/env -S deno run --allow-write
import { Hono } from "npm:hono";

Deno.addSignalListener("SIGINT", () => {
  console.log("Launch aborted!");
  Deno.exit(1);
});

const countDown = async (count: number) => {
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

const startServer = () => {
  const app = new Hono();
  Deno.serve({ port: 8080 }, app.fetch);
};

const main = async () => {
  if (Deno.args[0] === "--countdown" || "-cd") {
    const count = parseInt(Deno.args[1]);
    await countDown(count || 0);
  } else if (Deno.args[0] === "server") {
    startServer();
  } else {
    await countUp();
  }
};

await main();
