const fmt = {
    dateFormat(dateStamp: number, format: string = 'YYYY-mm-dd HH:MM:SS') {
        // 根据时间戳生成当前时间, 单位 毫秒
        let date = new Date(dateStamp)
        const opt = {
            "Y+": date.getFullYear().toString(),        // 年
            "m+": (date.getMonth() + 1).toString(),     // 月
            "d+": date.getDate().toString(),            // 日
            "H+": date.getHours().toString(),           // 时
            "M+": date.getMinutes().toString(),         // 分
            "S+": date.getSeconds().toString()          // 秒
            // 有其他格式化字符需求可以继续添加，必须转化成字符串
        };
        for (let k in opt) {
            // 正则匹配对应的 opt key
            let r = new RegExp("(" + k + ")").exec(format);
            // 如果有匹配成功项
            if (r) {
                format = format.replace(r[1], (r[1].length == 1) ? (opt[k as keyof typeof opt]) : (opt[k as keyof typeof opt].padStart(r[1].length, "0")))
            }
        }
        return format;
    },
    secondToTime(seconds: number) {
        const hours = Math.floor(seconds / 3600);
        const minutes = Math.floor((seconds % 3600) / 60);
        const remainingSeconds = Math.floor(seconds % 60);
        let timeStr = ''
        timeStr = hours < 10 ? `0${hours}` : `${hours}`;
        timeStr += minutes < 10 ? `:0${minutes}` : `:${minutes}`;
        timeStr += remainingSeconds < 10 ? `:0${remainingSeconds}` : `:${remainingSeconds}`;
        return timeStr;
    }
}

export {fmt}