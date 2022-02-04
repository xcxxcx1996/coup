import { ROLES, rolesMap } from "../../constants";
import { Box } from "@mui/material";
import Carousel from "react-material-ui-carousel";

export interface IUser {
    name: string;
    coins: number;
    roles: string[];
}

export interface UserCarouselProps {
    users: IUser[];
}

export interface UserInfoProps {
    user: IUser;
}

const users: IUser[] = [
    {
        name: "leo",
        coins: 2,
        roles: [rolesMap[ROLES.ASSASSIN], rolesMap[ROLES.QUEEN]],
    },
    {
        name: "john",
        coins: 2,
        roles: [rolesMap[ROLES.AMBASSADOR], rolesMap[ROLES.CAPTAIN]],
    },
    {
        name: "jack",
        coins: 2,
        roles: [rolesMap[ROLES.LORD], rolesMap[ROLES.ASSASSIN]],
    },
];

export function UserInfo(props: UserInfoProps) {
    const { user } = props;
    return (
        <Box
            sx={{
                p: 1,
                m: 1,
                mb: 2,
                display: "flex",
                flexDirection: "column",
                alignItems: "center",
                border: "1px solid",
                borderRadius: 2,
            }}
        >
            <div>玩家: {user.name}</div>
            <div>金币: {user.coins}</div>
            <div>角色: {user.roles.join(" | ")}</div>
        </Box>
    );
}

export function UserCarousel() {
    return (
        <Box
            sx={{
                mb: 1,
            }}
        >
            <Carousel
                autoPlay={false}
                indicators={false}
                navButtonsAlwaysVisible={true}
                animation={"slide"}
            >
                {users.map((user) => (
                    <UserInfo key={user.name} user={user} />
                ))}
            </Carousel>
        </Box>
    );
}
