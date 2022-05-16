import {
    Client,
    Intents,
    Interaction,
    Collection,
    MessageEmbed,
    GuildTextBasedChannel,
} from "discord.js";
import { REST } from "@discordjs/rest";
import { Routes } from "discord-api-types/v9";
import dotenv from "dotenv";
import path from "path";
import fs from "fs";
import cron from "node-cron";
import axios from "axios";

dotenv.config({
    path: process.env.ENV_PATH
        ? process.env.ENV_PATH
        : path.join(__dirname, ".env"),
});

const client = new Client({
    intents: [Intents.FLAGS.GUILDS, Intents.FLAGS.GUILD_MESSAGES],
});

const commands = new Collection();
const commandsData = [];
const commandsPath = path.join(__dirname, "commands");
const commandFiles = fs
    .readdirSync(commandsPath)
    .filter((file) => file.endsWith(".ts"));

for (const file of commandFiles) {
    const filePath = path.join(commandsPath, file);
    const command = require(filePath);
    commands.set(command.data.name, command);
    commandsData.push(command.data.toJSON());
}

client.once("ready", () => {
    console.log("the bot is ready");

    const rest = new REST({ version: "9" }).setToken(process.env.TOKEN);

    (async () => {
        try {
            if (process.env.ENV === "production") {
                console.log(
                    "[production] Started refreshing application (/) commands."
                );

                await rest.put(
                    // @ts-ignore
                    Routes.applicationCommands(client.user.id),
                    { body: commandsData }
                );

                console.log(
                    "[production] Successfully reloaded application (/) commands."
                );
            } else {
                console.log(
                    "[test] Started refreshing application (/) commands."
                );

                await rest.put(
                    // @ts-ignore
                    Routes.applicationGuildCommands(
                        client.user.id,
                        process.env.GUILD_ID
                    ),
                    { body: commandsData }
                );

                console.log(
                    "[test] Successfully reloaded application (/) commands."
                );
            }
        } catch (error) {
            console.error(error);
        }
    })();
});

client.on("interactionCreate", async (interaction: Interaction) => {
    if (!interaction.isCommand()) return;

    const command = commands.get(interaction.commandName);

    if (!command) return;

    try {
        // @ts-ignore
        await command.execute(interaction);
    } catch (error) {
        console.error(error);
        await interaction.reply({
            content: "There was an error while executing this command!",
            ephemeral: true,
        });
    }
});

client.login(process.env.TOKEN);

function newAlert(percentage: string, poolName: string): MessageEmbed {
    return new MessageEmbed()
        .setColor("RED")
        .setTitle("Alert!")
        .addField("\u200B", "\u200B")
        .addField("Pair address", poolName, true)
        .addField("\u200B", "\u200B", true)
        .addField("Percentage", percentage, true)
        .setAuthor({
            iconURL:
                "https://play-lh.googleusercontent.com/0bVs9-3xq573KI9u2hqZ86ARwltcoBv4RGOTI58Sw-xClAfl8dYdd9eYn2vf0D2HMA",
            name: "Alert bot",
        });
}

function pushMessage(
    client: Client,
    guild_id: string,
    channel_id: string,
    message?: string,
    embed?: MessageEmbed[]
) {
    const guild = client.guilds.cache.find((guild) => guild.id === guild_id);
    if (!guild) return;
    if (!guild.available) return;

    const channel = guild.channels.cache.find(
        (channel) => channel.id === channel_id
    );
    if (!channel) return;
    if (!channel.isText()) return;

    const text_channel: GuildTextBasedChannel = channel;

    text_channel.send({
        content: message,
        embeds: embed,
    });
}

let temp = new Map<string, number>();

// run alert check updated every 2 minute
cron.schedule("*/2 * * * *", () => {
    const tm = Math.round(new Date().getTime() / 1000);
    console.log("Now is " + tm.toString());

    interface Alert {
        address: string;
        high: number;
        low: number;
        change: number;
    }

    axios.get(process.env.API_URL + "/api/alert").then((response) => {
        let data: Alert[];
        data = response.data;

        data.map((alert) => {
            if (temp.has(alert.address)) {
                if (tm - temp.get(alert.address) > 3600) {
                }
            }
        });
    });
});
