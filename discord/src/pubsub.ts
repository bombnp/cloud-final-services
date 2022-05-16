import { Message, PubSub } from '@google-cloud/pubsub'
import { MessageEmbed } from 'discord.js'
import { pushMessage } from './client'

const TOPIC_PRICE_ALERT = 'price-alerts'
const TOPIC_PRICE_SUMMARY = 'price-summary'

const SUB_PRICE_ALERT = `${TOPIC_PRICE_ALERT}-${process.env.PUBSUB_SUBSCRIPTION_SUFFIX}`
const SUB_PRICE_SUMMARY = `${TOPIC_PRICE_SUMMARY}-${process.env.PUBSUB_SUBSCRIPTION_SUFFIX}`

interface PriceAlertMsg {
    serverId: string
    poolAddress: string
    channelId: string
    pairName: string
    change: number
    since: number
}

function onReceiveAlert(message: Message) {
    const alerts: PriceAlertMsg[] = JSON.parse(message.data.toString())
    for (const alert of alerts) {
        const embed = new MessageEmbed()
            .setColor('RED')
            .setTitle('Alert!')
            .addField('Pair', alert.pairName, true)
            .addField('\u200B', '\u200B', true)
            .addField('Address', alert.poolAddress, true)
            .addField('Change', `${(alert.change * 100).toFixed(2)}%`, true)
            .addField('\u200B', '\u200B', true)
            // TODO: format this nicely
            .addField('Since', new Date(alert.since * 1000).toISOString(), true)
            .setAuthor({
                iconURL:
                    'https://play-lh.googleusercontent.com/0bVs9-3xq573KI9u2hqZ86ARwltcoBv4RGOTI58Sw-xClAfl8dYdd9eYn2vf0D2HMA',
                name: 'Alert bot',
            })
        pushMessage(alert.serverId, alert.channelId, null, [embed])
    }
    message.ack()
}

const pubsub = new PubSub({
    projectId: process.env.PUBSUB_PROJECT_ID,
})

async function ensureSubscriptions() {
    const [subscriptions] = await pubsub.getSubscriptions()
    if (!subscriptions.find((sub) => sub.name.endsWith(SUB_PRICE_ALERT))) {
        await pubsub
            .topic(TOPIC_PRICE_ALERT)
            .createSubscription(SUB_PRICE_ALERT)
    }
    if (!subscriptions.find((sub) => sub.name.endsWith(SUB_PRICE_SUMMARY))) {
        await pubsub
            .topic(TOPIC_PRICE_ALERT)
            .createSubscription(SUB_PRICE_SUMMARY)
    }
}

export async function initPubSub() {
    await ensureSubscriptions()
    const priceAlertSub = pubsub
        .topic(TOPIC_PRICE_ALERT)
        .subscription(SUB_PRICE_ALERT)
    const priceSummarySub = pubsub
        .topic(TOPIC_PRICE_SUMMARY)
        .subscription(SUB_PRICE_SUMMARY)

    priceAlertSub.on('message', onReceiveAlert)

    console.log('subscribed to pubsub!')
}
