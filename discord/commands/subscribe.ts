import { SlashCommandBuilder } from "@discordjs/builders";
import { CommandInteraction } from "discord.js";

module.exports = {
    data: new SlashCommandBuilder()
        .setName("subscribe")
        .setDescription("subscribe to alert bot"),
    async execute(interaction: CommandInteraction) {
        await interaction.reply("I am still alive!");
    },
};
