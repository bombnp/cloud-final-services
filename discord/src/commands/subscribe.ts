import {
    SlashCommandBuilder,
    SlashCommandChannelOption,
    SlashCommandStringOption,
} from '@discordjs/builders'
import axios from 'axios'
import { ChannelType } from 'discord-api-types/v10'
import {
    CommandInteraction,
    GuildTextBasedChannel,
    MessageEmbed,
} from 'discord.js'
import dotenv from 'dotenv'

dotenv.config()

interface Pair {
    pool_address: string
    pool_name: string
    is_base_token0: string
}

async function getPairChoice(): Promise<Pair[]> {
    let pair_list: Pair[]
    await axios
        .get<Pair[]>(process.env.API_URL + '/api/pair')
        .then((response) => {
            pair_list = response.data
        })
    return pair_list
}

function SubscriptionEmbed(
    address: string,
    channel: GuildTextBasedChannel,
    successful = true
): MessageEmbed {
    return new MessageEmbed()
        .setColor(successful ? 'GREEN' : 'RED')
        .setTitle(`Subscription ${successful ? 'complete' : 'failed'}`)
        .addField('\u200B', '\u200B')
        .addField('Pair', address, true)
        .addField('\u200B', '\u200B', true)
        .addField('Alert channel', channel.toString(), true)
}

module.exports = {
    async init() {
        let command

        await getPairChoice().then((pair_list) => {
            command = new SlashCommandBuilder()
                .setName('subscribe')
                .setDescription('Subscribe to alert bot')
                .addStringOption((option: SlashCommandStringOption) => {
                    option = option
                        .setName('pair')
                        .setDescription('Pair that you wanna subscribe')
                        .setRequired(true)

                    pair_list.forEach((pair) => {
                        console.log(pair.pool_name + ': ' + pair.pool_address)
                        option = option.addChoices({
                            name: pair.pool_name,
                            value: pair.pool_address,
                        })
                    })

                    console.log(option.toJSON())

                    return option
                })
                .addChannelOption((option: SlashCommandChannelOption) => {
                    return option
                        .setName('channel')
                        .setDescription(
                            'Channel that you want for showing alert'
                        )
                        .setRequired(false)
                        .addChannelTypes(ChannelType.GuildText)
                })
        })

        return command
    },
    async execute(interaction: CommandInteraction) {
        const id = interaction.guildId
        const pair = interaction.options.getString('pair')
        const channel = interaction.options.getChannel('channel')
        let channel_target: GuildTextBasedChannel

        if (!channel) channel_target = interaction.channel
        else {
            channel_target = channel as GuildTextBasedChannel
        }

        try {
            await axios.post(process.env.API_URL + '/api/subscribe/alert', {
                server_id: id,
                pool: pair,
                channel_id: channel_target.id,
            })
            console.log(`[subscribe] subscribed ${id} for pair ${pair}`)

            const embed = SubscriptionEmbed(pair, channel_target, true)
            await interaction.reply({
                ephemeral: true,
                embeds: [embed],
            })
        } catch (err) {
            console.log(err.response.data)

            const embed = SubscriptionEmbed(pair, channel_target, false)
            await interaction.reply({
                ephemeral: true,
                embeds: [embed],
            })
        }
    },
}
