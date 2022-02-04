import { Container } from "@mui/material";
import { Header } from "../../components/in-game/Header";
import { UserInfos } from "../../components/in-game/UserInfos";

function InGame() {
    return (
        <Container maxWidth="xs">
            <Header />
            <UserInfos />
        </Container>
    );
}

export default InGame;
