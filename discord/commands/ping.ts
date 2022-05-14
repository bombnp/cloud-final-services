import { SlashCommandBuilder } from "@discordjs/builders";
import { CommandInteraction } from "discord.js";

module.exports = {
	data: new SlashCommandBuilder()
		.setName("ping")
		.setDescription("Test that the bot is still alive!"),
	async execute(interaction: CommandInteraction) {
		await interaction.reply({
			content: "I am still alive!",
			ephemeral: true,
		});
	},
};
