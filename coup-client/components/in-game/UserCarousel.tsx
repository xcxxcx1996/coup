import { Box } from "@mui/material";
import Carousel from "react-material-ui-carousel";
import { useContext } from "react";
import { gameContext } from "../../contexts/gameContext";

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
    const { users } = useContext(gameContext);
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
