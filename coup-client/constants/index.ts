export enum ROLES {
    UNROLE = 0,
    DIPLOMAT = 1,
    QUEEN = 2,
    CAPTAIN = 3,
    ASSASSIN = 4,
    BARON = 5,
}

export const rolesMap: { [role: string]: string } = {
    [ROLES.DIPLOMAT]: "大使",
    [ROLES.QUEEN]: "女王",
    [ROLES.BARON]: "男爵",
    [ROLES.ASSASSIN]: "刺客",
    [ROLES.CAPTAIN]: "队长",
    [ROLES.UNROLE]: "未知",
};
