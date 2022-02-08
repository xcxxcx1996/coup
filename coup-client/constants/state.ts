export const enum State {
    //undefined
    UNDEFINED = 0,
    //什么都做不了
    IDLE = 1,
    // 自己回合开始，可以拿1，2，3块钱，和施放技能
    START = 2,
    // 质疑
    QUESTION = 3,
    // 保护自己，包括防偷牌，看牌，和刺杀
    DENY_MONEY = 4,
    DENY_STEAL = 5,
    DENY_ASSASSIN = 6,
    DISCARD = 7,
    CHOOSE_CARD = 8,
    // DENY_CHECK=5;
}
