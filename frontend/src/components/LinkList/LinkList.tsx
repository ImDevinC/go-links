import { Box, Typography } from "@mui/joy";
import React, { useCallback, useEffect, useState } from "react";
import { LinkData, GetLinksResponse } from "../../services/api/links";
import { Link } from "../Link";

interface LinkListProps {
    getLinkFn: () => Promise<GetLinksResponse>
}

export const LinkList = (props: LinkListProps) => {
    const [links, setLinks] = useState<LinkData[]>([])

    const updateLinkList = useCallback(async () => {
        const links = await props.getLinkFn()
        setLinks(links.links)
    }, [props])

    useEffect(() => {
        updateLinkList()
    }, [updateLinkList])

    return (
        <Box sx={{ flex: 1, width: '100%', mx: 'auto' }}>
            {links.length > 0 ? links.map(link => <Link link={link} key={link.url} onLinkChanged={updateLinkList} />) : <Typography level="body-md">No links found</Typography>}
        </Box>



    )
}