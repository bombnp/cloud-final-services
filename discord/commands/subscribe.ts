import { SlashCommandBuilder } from "@discordjs/builders";
import { CommandInteraction } from "discord.js";

module.exports = {
	data: new SlashCommandBuilder(),
	async execute(interaction: CommandInteraction) {
		await interaction.reply("I am still alive!");
	},
};
