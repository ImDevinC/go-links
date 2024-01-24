import { ContentCopy, Edit, Visibility } from "@mui/icons-material";
import { Divider, Stack, Typography } from "@mui/joy";
import React from "react";
import { LinkData } from "../../services/api/links";

export interface LinkProps {
    link: LinkData
}

export const Link = (props: LinkProps) => {
    return (
        <Stack direction="column" spacing={2} my={2}>
            <Stack direction="row" spacing={2} my={2}>
                <Stack direction="column" spacing={1} flexGrow={1}>
                    <Typography level="body-md">go/{props.link.name}</Typography>
                    <Typography level="body-sm">{props.link.description}</Typography>
                </Stack>
                <Stack direction="row" spacing={2} sx={{ alignItems: 'center' }}>
                    <Typography level="body-md" startDecorator={<Visibility />}>{props.link.views}</Typography>
                    <Edit />
                    <ContentCopy />
                </Stack>
            </Stack>
            <Divider />
        </Stack>
    )
}