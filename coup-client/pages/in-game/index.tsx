import { Container } from "@mui/material";
import { Header } from "../../components/in-game/Header";
import { UserCarousel } from "../../components/in-game/UserCarousel";
import { GameHistory } from "../../components/in-game/GameHistory";
import { ControlPanel } from "../../components/in-game/ControlPanel";

function InGame() {
    return (
        <Container
            maxWidth="xs"
            sx={{
                minHeight: "100vh",
            }}
        >
            <Header />
            <UserCarousel />
            <GameHistory />
            <ControlPanel />
        </Container>
    );
}

export default InGame;
