import {
    Box,
    Button,
    Dialog,
    DialogTitle,
    Menu,
    MenuItem,
} from "@mui/material";
import React from "react";
import KeyboardArrowDownIcon from "@mui/icons-material/KeyboardArrowDown";
import { text } from "stream/consumers";
import { nakamaClient } from "../../utils/nakama";

export interface IMenuItem {
    text: string;
    onClick: () => void;
}

export interface MenuButtonProps {
    text: string;
    items: IMenuItem[];
}

const MenuButton = (props: MenuButtonProps) => {
    const { text, items } = props;
    const [anchorEl, setAnchorEl] = React.useState<null | HTMLElement>(null);
    const open = Boolean(anchorEl);
    const handleClick = (event: React.MouseEvent<HTMLButtonElement>) => {
        setAnchorEl(event.currentTarget);
    };
    const handleClose = () => {
        setAnchorEl(null);
    };
    return (
        <div>
            <Button
                sx={{ width: "140px", m: 1 }}
                variant="contained"
                onClick={handleClick}
                endIcon={<KeyboardArrowDownIcon />}
                disableElevation={true}
            >
                {text}
            </Button>
            <Menu
                id="basic-menu"
                anchorEl={anchorEl}
                open={open}
                onClose={handleClose}
                MenuListProps={{
                    "aria-labelledby": "basic-button",
                }}
            >
                {items.map((item) => (
                    <MenuItem key={item.text} onClick={item.onClick}>
                        {item.text}
                    </MenuItem>
                ))}
            </Menu>
        </div>
    );
};

export interface AbilityProps {
    open: boolean;
    handleClose: () => void;
}

export interface IAbility {
    text: string;
    ability: () => void;
}

export const AbilityDialog = (props: AbilityProps) => {
    const { open, handleClose } = props;

    const abilities: IAbility[] = [
        {
            text: "女王（防刺杀）",
            ability: async () => {
                await nakamaClient.denyKill();
            },
        },
        { text: "男爵（收3金币）", ability: () => {} },
        { text: "男爵（阻止收2金币）", ability: () => {} },
        { text: "大使（换牌）", ability: () => {} },
        { text: "大使（阻止偷金币）", ability: () => {} },
        { text: "队长（偷2金币）", ability: () => {} },
        { text: "队长（阻止偷金币）", ability: () => {} },
        { text: "刺客（刺杀）", ability: () => {} },
    ];

    return (
        <Dialog onClose={handleClose} open={open}>
            <DialogTitle sx={{ textAlign: "center" }}>选择卡牌</DialogTitle>
            <Box
                sx={{
                    display: "flex",
                    flexDirection: "column",
                    alignItems: "center",
                    p: 1,
                }}
            >
                {abilities.map((item) => (
                    <Button
                        key={item.text}
                        variant="contained"
                        sx={{ width: "250px", m: 1 }}
                        onClick={item.ability}
                    >
                        {item.text}
                    </Button>
                ))}
            </Box>
        </Dialog>
    );
};

export const ControlPanel = () => {
    const [open, setOpen] = React.useState(false);
    const handleClose = () => {
        setOpen(false);
    };
    const handleClickOpen = () => {
        setOpen(true);
    };
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
            <MenuButton
                items={[
                    { text: "1金币", onClick: () => {} },
                    { text: "2金币", onClick: () => {} },
                ]}
                text="获取金币"
            />
            <Button
                sx={{ width: "140px", m: 1 }}
                variant="contained"
                onClick={handleClickOpen}
            >
                使用技能
            </Button>
            <MenuButton
                items={[
                    { text: "质疑", onClick: () => {} },
                    { text: "不质疑", onClick: () => {} },
                ]}
                text="质疑/不质疑"
            />
            <Button sx={{ width: "140px", m: 1 }} variant="contained">
                政变
            </Button>
            <AbilityDialog open={open} handleClose={handleClose} />
        </Box>
    );
};
