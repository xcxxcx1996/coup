export enum ROLES {
    DIPLOMAT = 0,
    QUEEN = 1,
    CAPTAIN = 2,
    ASSASSIN = 3,
    BARON = 4,
}

export const rolesMap: { [role: string]: string } = {
    [ROLES.DIPLOMAT]: "大使",
    [ROLES.QUEEN]: "女王",
    [ROLES.BARON]: "男爵",
    [ROLES.ASSASSIN]: "刺客",
    [ROLES.CAPTAIN]: "队长",
};
