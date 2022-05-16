/**
 * Formats duration for display.
 * @param duration Duration in seconds
 * @returns Formatted duration
 */
export function formatDuration(duration: number): string {
    const seconds = duration % 60
    const minutes = Math.floor(duration / 60) % 60
    const hours = Math.floor(duration / 3600)
    if (hours == 0) {
        if (minutes == 0) {
            return `${seconds} seconds ago`
        }
        return `${minutes} minutes ago`
    }
    return `${hours} hours ago`
}
