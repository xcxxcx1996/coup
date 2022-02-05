import {
    Box,
    Button,
    Dialog,
    DialogTitle,
    Menu,
    MenuItem,
} from "@mui/material";
import React, { useContext, useState } from "react";
import KeyboardArrowDownIcon from "@mui/icons-material/KeyboardArrowDown";
import { gameContext } from "../../contexts/gameContext";

export interface IMenuItem {
    text: string;
    onClick: () => void;
}

export interface MenuButtonProps {
    text: string;
    items: IMenuItem[];
    btnWidth: number;
    menuItemWidth?: number;
}

const MenuButton = (props: MenuButtonProps) => {
    const { text, items, btnWidth, menuItemWidth = 140 } = props;
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
                sx={{ width: `${btnWidth}px`, m: 1 }}
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
                    <MenuItem
                        dense={true}
                        divider={true}
                        key={item.text}
                        onClick={item.onClick}
                        sx={{
                            width: `${menuItemWidth}px`,
                            display: "flex",
                            justifyContent: "center",
                        }}
                    >
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

export interface AbilityBtnProps {
    text: string;
    onClick: () => void;
}

export const AbilityBtn = ({ text, onClick }: AbilityBtnProps) => {
    return (
        <Button
            variant="contained"
            sx={{ width: "250px", m: 1 }}
            onClick={onClick}
        >
            {text}
        </Button>
    );
};

export const AbilityDialog = (props: AbilityProps) => {
    const { open, handleClose } = props;
    const { users } = useContext(gameContext);
    const targetUsers: IMenuItem[] = users.map((u) => ({
        text: u.name,
        onClick: null,
    }));
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
                <AbilityBtn text="女王（防刺杀）" onClick={null} />
                <AbilityBtn text="男爵（收3金币）" onClick={null} />
                <AbilityBtn text="男爵（阻止收2金币）" onClick={null} />
                <AbilityBtn text="大使（换牌）" onClick={null} />
                <AbilityBtn text="大使（阻止偷金币）" onClick={null} />
                <MenuButton
                    text={"队长（偷2金币）"}
                    items={targetUsers}
                    btnWidth={250}
                    menuItemWidth={250}
                />
                <AbilityBtn text="队长（阻止偷金币）" onClick={null} />
                <MenuButton
                    text={"刺客（刺杀）"}
                    items={targetUsers}
                    btnWidth={250}
                    menuItemWidth={250}
                />
            </Box>
        </Dialog>
    );
};

export const ControlPanel = () => {
    const [open, setOpen] = useState(false);
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
                btnWidth={140}
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
                btnWidth={140}
            />
            <Button sx={{ width: "140px", m: 1 }} variant="contained">
                政变
            </Button>
            <AbilityDialog open={open} handleClose={handleClose} />
        </Box>
    );
};
