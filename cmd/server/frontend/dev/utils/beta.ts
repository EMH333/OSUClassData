export const BetaEnabledKey = 'betaEnabled';

export function isBetaEnabled(): boolean {
    return localStorage.getItem(BetaEnabledKey) === 'true';
}
