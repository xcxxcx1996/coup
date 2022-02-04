import { Box, Button } from "@mui/material";

export const ControlPanel = () => {
    return (
        <Box
            sx={{
                display: "flex",
                flexDirection: "row",
                flexWrap: "wrap",
                justifyContent: "space-between",
                m: 1,
            }}
        >
            <Button sx={{ width: "120px", m: 1 }} variant="contained">
                获取金币
            </Button>
            <Button sx={{ width: "120px", m: 1 }} variant="contained">
                使用技能
            </Button>
            <Button sx={{ width: "120px", m: 1 }} variant="contained">
                质疑/不质疑
            </Button>
            <Button sx={{ width: "120px", m: 1 }} variant="contained">
                政变
            </Button>
        </Box>
    );
};
