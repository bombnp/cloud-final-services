import {
    EmbedBuilder,
    SlashCommandBuilder,
    SlashCommandChannelOption,
    SlashCommandStringOption,
} from "@discordjs/builders";
import { CommandInteraction, MessageEmbed } from "discord.js";

let command = new SlashCommandBuilder()
    .setName("subscribe")
    .setDescription("Subscribe to alert bot")
    .addStringOption((option: SlashCommandStringOption) => {
        option = option
            .setName("pair")
            .setDescription("Pair that you wanna subscribe")
            .setRequired(true)
            .addChoices({
                name: "BTC",
                value: "btc/es",
            });

        return option;
    })
    .addChannelOption((option: SlashCommandChannelOption) => {
        return option
            .setName("channel")
            .setDescription("Channel that you want for showing alert")
            .setRequired(false);
    });

module.exports = {
    data: command,
    async execute(interaction: CommandInteraction) {
        const pair = interaction.options.getString("pair");
        const channel = interaction.options.getChannel("channel", false);
        let channel_target;

        if (!channel) channel_target = interaction.channel;
        else channel_target = channel;

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
