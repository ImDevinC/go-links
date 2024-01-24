const baseUrl = process.env.NODE_ENV === "development" ? "http://localhost:8080" : "https://api.example.com";

export interface LinkData {
    url: string;
    name: string;
    description: string;
    views?: number
}

interface CreateLinkResponse {
    error?: string;
}

interface GetLinksResponse {
    links: LinkData[];
    error?: string;
}

export const createLink = async (link: LinkData): Promise<CreateLinkResponse> => {
    try {
        const response = await fetch(`${baseUrl}/${link.name}`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(link)
        });
        if (!response.ok) {
            let message = 'Failed to complete'
            if (response.body) {
                message = response.body.toString()
            }
            return { error: message }
        }
        return {}
    } catch (exception: any) {
        return { error: 'Failed to communicate with server, please try your request again' }
    }
}

export const getPopular = async (): Promise<GetLinksResponse> => {
    try {
        const response = await fetch(`${baseUrl}/api/popular`);
        if (!response.ok) {
            let message = 'Failed to complete'
            if (response.body) {
                message = response.body.toString()
            }
            return { links: [], error: message }
        }
        const links: LinkData[] = await response.json()
        return { links }
    } catch (exception: any) {
        return { links: [], error: 'Failed to communicate with server, please try your request again' }
    }
}