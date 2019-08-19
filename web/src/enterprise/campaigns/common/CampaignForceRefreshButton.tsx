import { LoadingSpinner } from '@sourcegraph/react-loading-spinner'
import SyncIcon from 'mdi-react/SyncIcon'
import React, { useCallback, useState } from 'react'
import { map, mapTo } from 'rxjs/operators'
import { NotificationType } from '../../../../../shared/src/api/client/services/notifications'
import { ExtensionsControllerNotificationProps } from '../../../../../shared/src/extensions/controller'
import { dataOrThrowErrors, gql } from '../../../../../shared/src/graphql/graphql'
import * as GQL from '../../../../../shared/src/graphql/schema'
import { mutateGraphQL } from '../../../backend/graphql'

const forceRefreshCampaign = (args: GQL.IForceRefreshCampaignOnMutationArguments): Promise<void> =>
    mutateGraphQL(
        gql`
            mutation ForceRefreshCampaign($campaign: ID!, $extensionData: CampaignExtensionData!) {
                forceRefreshCampaign(campaign: $campaign, extensionData: $extensionData) {
                    id
                }
            }
        `,
        args
    )
        .pipe(
            map(dataOrThrowErrors),
            mapTo(undefined)
        )
        .toPromise()

interface Props extends ExtensionsControllerNotificationProps {
    campaign: Pick<GQL.ICampaign, 'id'>
    onComplete?: () => void
    className?: string
    buttonClassName?: string
}

/**
 * A button that force-refreshes a campaign.
 */
export const CampaignForceRefreshButton: React.FunctionComponent<Props> = ({
    campaign,
    onComplete,
    className = '',
    buttonClassName = 'btn-link text-decoration-none',
    extensionsController,
}) => {
    const [isLoading, setIsLoading] = useState(false)
    const onClick = useCallback<React.FormEventHandler>(
        async e => {
            e.preventDefault()
            setIsLoading(true)
            try {
                await forceRefreshCampaign({
                    campaign: campaign.id,
                    // TODO!(sqs): support refreshing
                    extensionData: { rawDiagnostics: [], rawFileDiffs: [] },
                })
                setIsLoading(false)
                if (onComplete) {
                    onComplete()
                }
            } catch (err) {
                setIsLoading(false)
                extensionsController.services.notifications.showMessages.next({
                    message: `Error force-refreshing campaign: ${err.message}`,
                    type: NotificationType.Error,
                })
            }
        },
        [campaign.id, onComplete, extensionsController.services.notifications.showMessages]
    )
    return (
        <button type="button" disabled={isLoading} className={`btn ${buttonClassName} ${className}`} onClick={onClick}>
            {isLoading ? <LoadingSpinner className="icon-inline" /> : <SyncIcon className="icon-inline" />} Refresh
        </button>
    )
}
