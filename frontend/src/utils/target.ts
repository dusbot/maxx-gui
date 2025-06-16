interface PortRange {
    start: number;
    end?: number;
}

interface ParsedResult {
    ip: string;
    mask: number;
    ports: PortRange[];
    valid: boolean;
    error?: string;
}

function parseIpPortRange(input: string): ParsedResult {
    const basicRegex = /^(\d{1,3}(\.\d{1,3}){3})\/(\d{1,2}):$$(.+)$$$/;
    const basicMatch = input.match(basicRegex);

    if (!basicMatch) {
        return { ip: '', mask: 0, ports: [], valid: false, error: 'Unknown format' };
    }

    const [, ip, , maskStr, portsStr] = basicMatch;
    const mask = parseInt(maskStr, 10);

    if (!isValidIp(ip)) {
        return { ip, mask, ports: [], valid: false, error: 'IP address is invalid' };
    }

    if (mask < 1 || mask > 32) {
        return { ip, mask, ports: [], valid: false, error: 'Subnet mask must be between 1 and 32' };
    }

    const ports: PortRange[] = [];
    const portParts = portsStr.split('|');

    for (const part of portParts) {
        if (part === '') {
            return { ip, mask, ports: [], valid: false, error: 'Port range cannot be empty' };
        }

        const rangeMatch = part.match(/^(\d+)(?:-(\d+))?$/);
        if (!rangeMatch) {
            return { ip, mask, ports: [], valid: false, error: `Port range format is invalid: ${part}` };
        }

        const [, startStr, endStr] = rangeMatch;
        const start = parseInt(startStr, 10);
        const end = endStr ? parseInt(endStr, 10) : undefined;

        if (start < 1 || start > 65535) {
            return { ip, mask, ports: [], valid: false, error: `Port number out of range: ${start}` };
        }

        if (end && (end < 1 || end > 65535)) {
            return { ip, mask, ports: [], valid: false, error: `Port number out of range: ${end}` };
        }

        if (end && end < start) {
            return { ip, mask, ports: [], valid: false, error: `Port range is invalid: ${start}-${end}` };
        }

        ports.push(end ? { start, end } : { start });
    }

    return { ip, mask, ports, valid: true };
}

function isValidIp(ip: string): boolean {
    const parts = ip.split('.').map(Number);
    return parts.length === 4 &&
        parts.every(part => part >= 0 && part <= 255) &&
        !parts.some(part => isNaN(part));
}

export {
    parseIpPortRange,
    isValidIp
}

