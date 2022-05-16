import {
    SlashCommandBuilder,
    SlashCommandChannelOption,
    SlashCommandStringOption,
} from "@discordjs/builders";
import axios from "axios";
import { ChannelType } from "discord-api-types/v10";
import { CommandInteraction, MessageEmbed } from "discord.js";
import dotenv from "dotenv";

dotenv.config();

interface Pair {
    pool_address: string;
    pool_name: string;
    is_base_token0: string;
}

async function getPairChoice(): Promise<Pair[]> {
    let pair_list: Pair[];
    await axios
        .get<Pair[]>(process.env.API_URL + "/api/pair")
        .then((response) => {
            pair_list = response.data;
        });
    return pair_list;
}

let command = new SlashCommandBuilder()
    .setName("subscribe")
    .setDescription("Subscribe to alert bot")
    .addStringOption((option: SlashCommandStringOption) => {
        option = option
            .setName("pair")
            .setDescription("Pair that you wanna subscribe")
            .setRequired(true);

        getPairChoice().then((pair_list) => {
            pair_list.map((pair) => {
                console.log(pair.pool_name + ": " + pair.pool_address);
                option.setChoices({
                    name: pair.pool_name,
                    value: pair.pool_address,
                });
            });
        });

        return option;
    })
    .addChannelOption((option: SlashCommandChannelOption) => {
        return option
            .setName("channel")
            .setDescription("Channel that you want for showing alert")
            .setRequired(false)
            .addChannelTypes(ChannelType.GuildText);
    });

module.exports = {
    data: command,
    async execute(interaction: CommandInteraction) {
        const id = interaction.guildId;
        const pair = interaction.options.getString("pair");
        const channel = interaction.options.getChannel("channel");
        let channel_target;

        if (!channel) channel_target = interaction.channel;
        else {
            channel_target = channel;
        }

        axios
            .post(process.env.API_URL + "/api/subscribe/alert", {
                server_id: id,
                pool: pair,
                channel: channel_target.id,
            })
            .catch((err) => {
                console.log(err.data);
            });

        const embed: MessageEmbed = new MessageEmbed()
            .setColor("GREEN")
            .setTitle("Subscribe complete")
            .setThumbnail(
                "https://play-lh.googleusercontent.com/0bVs9-3xq573KI9u2hqZ86ARwltcoBv4RGOTI58Sw-xClAfl8dYdd9eYn2vf0D2HMA"
            )
            .addField("\u200B", "\u200B")
            .addField("Pair address", pair, true)
            .addField("\u200B", "\u200B", true)
            .addField("Alert Channel", channel_target.toString(), true)
            .setAuthor({
                iconURL:
                    "https://play-lh.googleusercontent.com/0bVs9-3xq573KI9u2hqZ86ARwltcoBv4RGOTI58Sw-xClAfl8dYdd9eYn2vf0D2HMA",
                name: "Alert bot",
            });

        await interaction.reply({
            ephemeral: true,
            embeds: [embed],
        });
    },
};
