export function formateDate(dateString: string): string {
    const d = new Date(dateString);

    if (Number.isNaN(d.getTime())) {
        return dateString
    }

    return d.toLocaleString('ru-RU', {
        day: "2-digit",
        month: "2-digit",
        year: "numeric",
        hour: "2-digit",
        minute: "2-digit",
        second: "2-digit"
    })
}