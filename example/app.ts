import discord from "@splitscript.js/discord";

discord.listen("123");
let iter = 0;
setInterval(() => {
  console.log("avbc");
  iter++;
  if (iter === 10) {
    throw "abc";
  }
}, 1000);
