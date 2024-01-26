import { Box, Typography } from "@mui/joy";
import React, { useEffect, useState } from "react";
import { LinkData, GetLinksResponse } from "../../services/api/links";
import { Link } from "../Link";

interface LinkListProps {
    getLinkFn: () => Promise<GetLinksResponse>
}

export const LinkList = (props: LinkListProps) => {
    const [links, setLinks] = useState<LinkData[]>([])

    useEffect(() => {
        const fetchLinks = async () => {
            const links = await props.getLinkFn()
            setLinks(links.links)
        }
        fetchLinks()
    }, [props])

    return (
        <Box sx={{ flex: 1, width: '100%', mx: 'auto' }}>
            {links.length > 0 ? links.map(link => <Link link={link} key={link.url} />) : <Typography level="body-md">No links found</Typography>}
        </Box>



    )
}