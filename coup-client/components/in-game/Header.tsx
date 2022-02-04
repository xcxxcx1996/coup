import { Box } from "@mui/system";

export function Header() {
    return (
        <Box
            sx={{
                p: 1,
                m: 1,
                display: "flex",
                justifyContent: "space-between",
                color: "red",
            }}
        >
            <div>John的回合</div>
            <div>回合时间: 33s</div>
        </Box>
    );
}
