import { DarkMode, LightMode } from "@mui/icons-material";
import { Switch } from "@mui/joy";
import { useColorScheme } from "@mui/joy";
import React, { useState } from "react";

export const ThemeToggle = () => {
    const { mode, setMode } = useColorScheme();
    const [thumb, setThumb] = useState(mode === 'dark' ? <DarkMode /> : <LightMode />)

    const toggleMode = () => {
        setThumb(mode === 'dark' ? <LightMode /> : <DarkMode />)
        setMode(mode === 'dark' ? 'light' : 'dark')
    }

    return (
        <Switch
            size="lg"
            slotProps={{
                input: { 'arial-label': 'Dark mode' },
                thumb: {
                    children: thumb
                }
            }}
            checked={mode === 'light'}
            sx={{
                '--Switch-thumbSize': '16px',
                alignSelf: 'flex-end',
            }}
            onChange={toggleMode}
        />
    )
}