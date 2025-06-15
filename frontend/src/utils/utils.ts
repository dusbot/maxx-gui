export const genHash = (...args: string[]): string => {
    let combinedHash = 0;
    args.forEach(arg => {
        for (let i = 0; i < arg.length; i++) {
            const char = arg.charCodeAt(i);
            combinedHash = ((combinedHash << 5) - combinedHash) + char;
            combinedHash = combinedHash & combinedHash;
        }
    });
    return combinedHash.toString();
}

export function addUniqueItem<T extends Record<string, any>>(
    array: T[],
    newItem: T,
    keys: (keyof T)[]
): void {
    const isDuplicate = array.some((item) =>
        keys.every((key) => item[key] === newItem[key])
    );

    if (!isDuplicate) {
        array.push(newItem);
    }
}