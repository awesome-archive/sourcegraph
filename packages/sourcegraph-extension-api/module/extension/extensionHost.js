import { Subscription } from 'rxjs';
import { createProxy, handleRequests } from '../common/proxy';
import { createConnection } from '../protocol/jsonrpc2/connection';
import { createWebWorkerMessageTransports } from '../protocol/jsonrpc2/transports/webWorker';
import { ExtCommands } from './api/commands';
import { ExtConfiguration } from './api/configuration';
import { ExtContext } from './api/context';
import { ExtDocuments } from './api/documents';
import { ExtLanguageFeatures } from './api/languageFeatures';
import { ExtSearch } from './api/search';
import { ExtViews } from './api/views';
import { ExtWindows } from './api/windows';
import { Location } from './types/location';
import { Position } from './types/position';
import { Range } from './types/range';
import { Selection } from './types/selection';
import { URI } from './types/uri';
const consoleLogger = {
    error(message) {
        console.error(message);
    },
    warn(message) {
        console.warn(message);
    },
    info(message) {
        console.info(message);
    },
    log(message) {
        console.log(message);
    },
};
/**
 * Creates the Sourcegraph extension host and the extension API handle (which extensions access with `import
 * sourcegraph from 'sourcegraph'`).
 *
 * @param initData The information to initialize this extension host.
 * @param transports The message reader and writer to use for communication with the client. Defaults to
 *                   communicating using self.postMessage and MessageEvents with the parent (assuming that it is
 *                   called in a Web Worker).
 * @return The extension API.
 */
export function createExtensionHost(initData, transports = createWebWorkerMessageTransports()) {
    const connection = createConnection(transports, consoleLogger);
    connection.listen();
    return createExtensionHandle(initData, connection);
}
function createExtensionHandle(initData, connection) {
    const subscription = new Subscription();
    subscription.add(connection);
    // For debugging/tests.
    const sync = () => connection.sendRequest('ping');
    connection.onRequest('ping', () => 'pong');
    const proxy = (prefix) => createProxy((name, args) => connection.sendRequest(`${prefix}/${name}`, args));
    const context = new ExtContext(proxy('context'));
    handleRequests(connection, 'context', context);
    const documents = new ExtDocuments(sync);
    handleRequests(connection, 'documents', documents);
    const windows = new ExtWindows(proxy('windows'), proxy('codeEditor'), documents);
    handleRequests(connection, 'windows', windows);
    const views = new ExtViews(proxy('views'));
    handleRequests(connection, 'views', views);
    const configuration = new ExtConfiguration(proxy('configuration'));
    handleRequests(connection, 'configuration', configuration);
    const languageFeatures = new ExtLanguageFeatures(proxy('languageFeatures'), documents);
    handleRequests(connection, 'languageFeatures', languageFeatures);
    const search = new ExtSearch(proxy('search'));
    handleRequests(connection, 'search', search);
    const commands = new ExtCommands(proxy('commands'));
    handleRequests(connection, 'commands', commands);
    return {
        URI,
        Position,
        Range,
        Selection,
        Location,
        MarkupKind: {
            PlainText: "plaintext" /* PlainText */,
            Markdown: "markdown" /* Markdown */,
        },
        app: {
            get activeWindow() {
                return windows.getActive();
            },
            get windows() {
                return windows.getAll();
            },
            createPanelView: id => views.createPanelView(id),
        },
        workspace: {
            get textDocuments() {
                return documents.getAll();
            },
            onDidOpenTextDocument: documents.onDidOpenTextDocument,
        },
        configuration: {
            get: () => configuration.get(),
            subscribe: next => configuration.subscribe(next),
        },
        languages: {
            registerHoverProvider: (selector, provider) => languageFeatures.registerHoverProvider(selector, provider),
            registerDefinitionProvider: (selector, provider) => languageFeatures.registerDefinitionProvider(selector, provider),
            registerTypeDefinitionProvider: (selector, provider) => languageFeatures.registerTypeDefinitionProvider(selector, provider),
            registerImplementationProvider: (selector, provider) => languageFeatures.registerImplementationProvider(selector, provider),
            registerReferenceProvider: (selector, provider) => languageFeatures.registerReferenceProvider(selector, provider),
        },
        search: {
            registerQueryTransformer: provider => search.registerQueryTransformer(provider),
        },
        commands: {
            registerCommand: (command, callback) => commands.registerCommand({ command, callback }),
            executeCommand: (command, ...args) => commands.executeCommand(command, args),
        },
        internal: {
            sync,
            updateContext: updates => context.updateContext(updates),
            sourcegraphURL: new URI(initData.sourcegraphURL),
            clientApplication: initData.clientApplication,
        },
    };
}
//# sourceMappingURL=data:application/json;base64,eyJ2ZXJzaW9uIjozLCJmaWxlIjoiZXh0ZW5zaW9uSG9zdC5qcyIsInNvdXJjZVJvb3QiOiJzcmMvIiwic291cmNlcyI6WyJleHRlbnNpb24vZXh0ZW5zaW9uSG9zdC50cyJdLCJuYW1lcyI6W10sIm1hcHBpbmdzIjoiQUFBQSxPQUFPLEVBQUUsWUFBWSxFQUFFLE1BQU0sTUFBTSxDQUFBO0FBRW5DLE9BQU8sRUFBRSxXQUFXLEVBQUUsY0FBYyxFQUFFLE1BQU0saUJBQWlCLENBQUE7QUFDN0QsT0FBTyxFQUFjLGdCQUFnQixFQUE2QixNQUFNLGlDQUFpQyxDQUFBO0FBQ3pHLE9BQU8sRUFBRSxnQ0FBZ0MsRUFBRSxNQUFNLDJDQUEyQyxDQUFBO0FBQzVGLE9BQU8sRUFBRSxXQUFXLEVBQUUsTUFBTSxnQkFBZ0IsQ0FBQTtBQUM1QyxPQUFPLEVBQUUsZ0JBQWdCLEVBQUUsTUFBTSxxQkFBcUIsQ0FBQTtBQUN0RCxPQUFPLEVBQUUsVUFBVSxFQUFFLE1BQU0sZUFBZSxDQUFBO0FBQzFDLE9BQU8sRUFBRSxZQUFZLEVBQUUsTUFBTSxpQkFBaUIsQ0FBQTtBQUM5QyxPQUFPLEVBQUUsbUJBQW1CLEVBQUUsTUFBTSx3QkFBd0IsQ0FBQTtBQUM1RCxPQUFPLEVBQUUsU0FBUyxFQUFFLE1BQU0sY0FBYyxDQUFBO0FBQ3hDLE9BQU8sRUFBRSxRQUFRLEVBQUUsTUFBTSxhQUFhLENBQUE7QUFDdEMsT0FBTyxFQUFFLFVBQVUsRUFBRSxNQUFNLGVBQWUsQ0FBQTtBQUMxQyxPQUFPLEVBQUUsUUFBUSxFQUFFLE1BQU0sa0JBQWtCLENBQUE7QUFDM0MsT0FBTyxFQUFFLFFBQVEsRUFBRSxNQUFNLGtCQUFrQixDQUFBO0FBQzNDLE9BQU8sRUFBRSxLQUFLLEVBQUUsTUFBTSxlQUFlLENBQUE7QUFDckMsT0FBTyxFQUFFLFNBQVMsRUFBRSxNQUFNLG1CQUFtQixDQUFBO0FBQzdDLE9BQU8sRUFBRSxHQUFHLEVBQUUsTUFBTSxhQUFhLENBQUE7QUFFakMsTUFBTSxhQUFhLEdBQVc7SUFDMUIsS0FBSyxDQUFDLE9BQWU7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QztZQUNuRSw4QkFBOEIsRUFBRSxDQUFDLFFBQVEsRUFBRSxRQUFRLEVBQUUsRUFBRSxDQUNuRCxnQkFBZ0IsQ0FBQyw4QkFBOEIsQ0FBQyxRQUFRLEVBQUUsUUFBUSxDQUFDO1lBQ3ZFLDhCQUE4QixFQUFFLENBQUMsUUFBUSxFQUFFLFFBQVEsRUFBRSxFQUFFLENBQ25ELGdCQUFnQixDQUFDLDhCQUE4QixDQUFDLFFBQVEsRUFBRSxRQUFRLENBQUM7WUFDdkUseUJBQXlCLEVBQUUsQ0FBQyxRQUFRLEVBQUUsUUFBUSxFQUFFLEVBQUUsQ0FDOUMsZ0JBQWdCLENBQUMseUJBQXlCLENBQUMsUUFBUSxFQUFFLFFBQVEsQ0FBQztTQUNyRTtRQUVELE1BQU0sRUFBRTtZQUNKLHdCQUF3QixFQUFFLFFBQVEsQ0FBQyxFQUFFLENBQUMsTUFBTSxDQUFDLHdCQUF3QixDQUFDLFFBQVEsQ0FBQztTQUNsRjtRQUVELFFBQVEsRUFBRTtZQUNOLGVBQWUsRUFBRSxDQUFDLE9BQU8sRUFBRSxRQUFRLEVBQUUsRUFBRSxDQUFDLFFBQVEsQ0FBQyxlQUFlLENBQUMsRUFBRSxPQUFPLEVBQUUsUUFBUSxFQUFFLENBQUM7WUFDdkYsY0FBYyxFQUFFLENBQUMsT0FBTyxFQUFFLEdBQUcsSUFBSSxFQUFFLEVBQUUsQ0FBQyxRQUFRLENBQUMsY0FBYyxDQUFDLE9BQU8sRUFBRSxJQUFJLENBQUM7U0FDL0U7UUFFRCxRQUFRLEVBQUU7WUFDTixJQUFJO1lBQ0osYUFBYSxFQUFFLE9BQU8sQ0FBQyxFQUFFLENBQUMsT0FBTyxDQUFDLGFBQWEsQ0FBQyxPQUFPLENBQUM7WUFDeEQsY0FBYyxFQUFFLElBQUksR0FBRyxDQUFDLFFBQVEsQ0FBQyxjQUFjLENBQUM7WUFDaEQsaUJBQWlCLEVBQUUsUUFBUSxDQUFDLGlCQUFpQjtTQUNoRDtLQUNKLENBQUE7QUFDTCxDQUFDIiwic291cmNlc0NvbnRlbnQiOlsiaW1wb3J0IHsgU3Vic2NyaXB0aW9uIH0gZnJvbSAncnhqcydcbmltcG9ydCAqIGFzIHNvdXJjZWdyYXBoIGZyb20gJ3NvdXJjZWdyYXBoJ1xuaW1wb3J0IHsgY3JlYXRlUHJveHksIGhhbmRsZVJlcXVlc3RzIH0gZnJvbSAnLi4vY29tbW9uL3Byb3h5J1xuaW1wb3J0IHsgQ29ubmVjdGlvbiwgY3JlYXRlQ29ubmVjdGlvbiwgTG9nZ2VyLCBNZXNzYWdlVHJhbnNwb3J0cyB9IGZyb20gJy4uL3Byb3RvY29sL2pzb25ycGMyL2Nvbm5lY3Rpb24nXG5pbXBvcnQgeyBjcmVhdGVXZWJXb3JrZXJNZXNzYWdlVHJhbnNwb3J0cyB9IGZyb20gJy4uL3Byb3RvY29sL2pzb25ycGMyL3RyYW5zcG9ydHMvd2ViV29ya2VyJ1xuaW1wb3J0IHsgRXh0Q29tbWFuZHMgfSBmcm9tICcuL2FwaS9jb21tYW5kcydcbmltcG9ydCB7IEV4dENvbmZpZ3VyYXRpb24gfSBmcm9tICcuL2FwaS9jb25maWd1cmF0aW9uJ1xuaW1wb3J0IHsgRXh0Q29udGV4dCB9IGZyb20gJy4vYXBpL2NvbnRleHQnXG5pbXBvcnQgeyBFeHREb2N1bWVudHMgfSBmcm9tICcuL2FwaS9kb2N1bWVudHMnXG5pbXBvcnQgeyBFeHRMYW5ndWFnZUZlYXR1cmVzIH0gZnJvbSAnLi9hcGkvbGFuZ3VhZ2VGZWF0dXJlcydcbmltcG9ydCB7IEV4dFNlYXJjaCB9IGZyb20gJy4vYXBpL3NlYXJjaCdcbmltcG9ydCB7IEV4dFZpZXdzIH0gZnJvbSAnLi9hcGkvdmlld3MnXG5pbXBvcnQgeyBFeHRXaW5kb3dzIH0gZnJvbSAnLi9hcGkvd2luZG93cydcbmltcG9ydCB7IExvY2F0aW9uIH0gZnJvbSAnLi90eXBlcy9sb2NhdGlvbidcbmltcG9ydCB7IFBvc2l0aW9uIH0gZnJvbSAnLi90eXBlcy9wb3NpdGlvbidcbmltcG9ydCB7IFJhbmdlIH0gZnJvbSAnLi90eXBlcy9yYW5nZSdcbmltcG9ydCB7IFNlbGVjdGlvbiB9IGZyb20gJy4vdHlwZXMvc2VsZWN0aW9uJ1xuaW1wb3J0IHsgVVJJIH0gZnJvbSAnLi90eXBlcy91cmknXG5cbmNvbnN0IGNvbnNvbGVMb2dnZXI6IExvZ2dlciA9IHtcbiAgICBlcnJvcihtZXNzYWdlOiBzdHJpbmcpOiB2b2lkIHtcbiAgICAgICAgY29uc29sZS5lcnJvcihtZXNzYWdlKVxuICAgIH0sXG4gICAgd2FybihtZXNzYWdlOiBzdHJpbmcpOiB2b2lkIHtcbiAgICAgICAgY29uc29sZS53YXJuKG1lc3NhZ2UpXG4gICAgfSxcbiAgICBpbmZvKG1lc3NhZ2U6IHN0cmluZyk6IHZvaWQge1xuICAgICAgICBjb25zb2xlLmluZm8obWVzc2FnZSlcbiAgICB9LFxuICAgIGxvZyhtZXNzYWdlOiBzdHJpbmcpOiB2b2lkIHtcbiAgICAgICAgY29uc29sZS5sb2cobWVzc2FnZSlcbiAgICB9LFxufVxuXG4vKipcbiAqIFJlcXVpcmVkIGluZm9ybWF0aW9uIHdoZW4gaW5pdGlhbGl6aW5nIGFuIGV4dGVuc2lvbiBob3N0LlxuICovXG5leHBvcnQgaW50ZXJmYWNlIEluaXREYXRhIHtcbiAgICAvKiogVGhlIFVSTCB0byB0aGUgSmF2YVNjcmlwdCBzb3VyY2UgZmlsZSAodGhhdCBleHBvcnRzIGFuIGBhY3RpdmF0ZWAgZnVuY3Rpb24pIGZvciB0aGUgZXh0ZW5zaW9uLiAqL1xuICAgIGJ1bmRsZVVSTDogc3RyaW5nXG5cbiAgICAvKiogQHNlZSB7QGxpbmsgbW9kdWxlOnNvdXJjZWdyYXBoLmludGVybmFsLnNvdXJjZWdyYXBoVVJMfSAqL1xuICAgIHNvdXJjZWdyYXBoVVJMOiBzdHJpbmdcblxuICAgIC8qKiBAc2VlIHtAbGluayBtb2R1bGU6c291cmNlZ3JhcGguaW50ZXJuYWwuY2xpZW50QXBwbGljYXRpb259ICovXG4gICAgY2xpZW50QXBwbGljYXRpb246ICdzb3VyY2VncmFwaCcgfCAnb3RoZXInXG59XG5cbi8qKlxuICogQ3JlYXRlcyB0aGUgU291cmNlZ3JhcGggZXh0ZW5zaW9uIGhvc3QgYW5kIHRoZSBleHRlbnNpb24gQVBJIGhhbmRsZSAod2hpY2ggZXh0ZW5zaW9ucyBhY2Nlc3Mgd2l0aCBgaW1wb3J0XG4gKiBzb3VyY2VncmFwaCBmcm9tICdzb3VyY2VncmFwaCdgKS5cbiAqXG4gKiBAcGFyYW0gaW5pdERhdGEgVGhlIGluZm9ybWF0aW9uIHRvIGluaXRpYWxpemUgdGhpcyBleHRlbnNpb24gaG9zdC5cbiAqIEBwYXJhbSB0cmFuc3BvcnRzIFRoZSBtZXNzYWdlIHJlYWRlciBhbmQgd3JpdGVyIHRvIHVzZSBmb3IgY29tbXVuaWNhdGlvbiB3aXRoIHRoZSBjbGllbnQuIERlZmF1bHRzIHRvXG4gKiAgICAgICAgICAgICAgICAgICBjb21tdW5pY2F0aW5nIHVzaW5nIHNlbGYucG9zdE1lc3NhZ2UgYW5kIE1lc3NhZ2VFdmVudHMgd2l0aCB0aGUgcGFyZW50IChhc3N1bWluZyB0aGF0IGl0IGlzXG4gKiAgICAgICAgICAgICAgICAgICBjYWxsZWQgaW4gYSBXZWIgV29ya2VyKS5cbiAqIEByZXR1cm4gVGhlIGV4dGVuc2lvbiBBUEkuXG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBjcmVhdGVFeHRlbnNpb25Ib3N0KFxuICAgIGluaXREYXRhOiBJbml0RGF0YSxcbiAgICB0cmFuc3BvcnRzOiBNZXNzYWdlVHJhbnNwb3J0cyA9IGNyZWF0ZVdlYldvcmtlck1lc3NhZ2VUcmFuc3BvcnRzKClcbik6IHR5cGVvZiBzb3VyY2VncmFwaCB7XG4gICAgY29uc3QgY29ubmVjdGlvbiA9IGNyZWF0ZUNvbm5lY3Rpb24odHJhbnNwb3J0cywgY29uc29sZUxvZ2dlcilcbiAgICBjb25uZWN0aW9uLmxpc3RlbigpXG4gICAgcmV0dXJuIGNyZWF0ZUV4dGVuc2lvbkhhbmRsZShpbml0RGF0YSwgY29ubmVjdGlvbilcbn1cblxuZnVuY3Rpb24gY3JlYXRlRXh0ZW5zaW9uSGFuZGxlKGluaXREYXRhOiBJbml0RGF0YSwgY29ubmVjdGlvbjogQ29ubmVjdGlvbik6IHR5cGVvZiBzb3VyY2VncmFwaCB7XG4gICAgY29uc3Qgc3Vic2NyaXB0aW9uID0gbmV3IFN1YnNjcmlwdGlvbigpXG4gICAgc3Vic2NyaXB0aW9uLmFkZChjb25uZWN0aW9uKVxuXG4gICAgLy8gRm9yIGRlYnVnZ2luZy90ZXN0cy5cbiAgICBjb25zdCBzeW5jID0gKCkgPT4gY29ubmVjdGlvbi5zZW5kUmVxdWVzdDx2b2lkPigncGluZycpXG4gICAgY29ubmVjdGlvbi5vblJlcXVlc3QoJ3BpbmcnLCAoKSA9PiAncG9uZycpXG5cbiAgICBjb25zdCBwcm94eSA9IChwcmVmaXg6IHN0cmluZykgPT4gY3JlYXRlUHJveHkoKG5hbWUsIGFyZ3MpID0+IGNvbm5lY3Rpb24uc2VuZFJlcXVlc3QoYCR7cHJlZml4fS8ke25hbWV9YCwgYXJncykpXG5cbiAgICBjb25zdCBjb250ZXh0ID0gbmV3IEV4dENvbnRleHQocHJveHkoJ2NvbnRleHQnKSlcbiAgICBoYW5kbGVSZXF1ZXN0cyhjb25uZWN0aW9uLCAnY29udGV4dCcsIGNvbnRleHQpXG5cbiAgICBjb25zdCBkb2N1bWVudHMgPSBuZXcgRXh0RG9jdW1lbnRzKHN5bmMpXG4gICAgaGFuZGxlUmVxdWVzdHMoY29ubmVjdGlvbiwgJ2RvY3VtZW50cycsIGRvY3VtZW50cylcblxuICAgIGNvbnN0IHdpbmRvd3MgPSBuZXcgRXh0V2luZG93cyhwcm94eSgnd2luZG93cycpLCBwcm94eSgnY29kZUVkaXRvcicpLCBkb2N1bWVudHMpXG4gICAgaGFuZGxlUmVxdWVzdHMoY29ubmVjdGlvbiwgJ3dpbmRvd3MnLCB3aW5kb3dzKVxuXG4gICAgY29uc3Qgdmlld3MgPSBuZXcgRXh0Vmlld3MocHJveHkoJ3ZpZXdzJykpXG4gICAgaGFuZGxlUmVxdWVzdHMoY29ubmVjdGlvbiwgJ3ZpZXdzJywgdmlld3MpXG5cbiAgICBjb25zdCBjb25maWd1cmF0aW9uID0gbmV3IEV4dENvbmZpZ3VyYXRpb248YW55Pihwcm94eSgnY29uZmlndXJhdGlvbicpKVxuICAgIGhhbmRsZVJlcXVlc3RzKGNvbm5lY3Rpb24sICdjb25maWd1cmF0aW9uJywgY29uZmlndXJhdGlvbilcblxuICAgIGNvbnN0IGxhbmd1YWdlRmVhdHVyZXMgPSBuZXcgRXh0TGFuZ3VhZ2VGZWF0dXJlcyhwcm94eSgnbGFuZ3VhZ2VGZWF0dXJlcycpLCBkb2N1bWVudHMpXG4gICAgaGFuZGxlUmVxdWVzdHMoY29ubmVjdGlvbiwgJ2xhbmd1YWdlRmVhdHVyZXMnLCBsYW5ndWFnZUZlYXR1cmVzKVxuXG4gICAgY29uc3Qgc2VhcmNoID0gbmV3IEV4dFNlYXJjaChwcm94eSgnc2VhcmNoJykpXG4gICAgaGFuZGxlUmVxdWVzdHMoY29ubmVjdGlvbiwgJ3NlYXJjaCcsIHNlYXJjaClcblxuICAgIGNvbnN0IGNvbW1hbmRzID0gbmV3IEV4dENvbW1hbmRzKHByb3h5KCdjb21tYW5kcycpKVxuICAgIGhhbmRsZVJlcXVlc3RzKGNvbm5lY3Rpb24sICdjb21tYW5kcycsIGNvbW1hbmRzKVxuXG4gICAgcmV0dXJuIHtcbiAgICAgICAgVVJJLFxuICAgICAgICBQb3NpdGlvbixcbiAgICAgICAgUmFuZ2UsXG4gICAgICAgIFNlbGVjdGlvbixcbiAgICAgICAgTG9jYXRpb24sXG4gICAgICAgIE1hcmt1cEtpbmQ6IHtcbiAgICAgICAgICAgIFBsYWluVGV4dDogc291cmNlZ3JhcGguTWFya3VwS2luZC5QbGFpblRleHQsXG4gICAgICAgICAgICBNYXJrZG93bjogc291cmNlZ3JhcGguTWFya3VwS2luZC5NYXJrZG93bixcbiAgICAgICAgfSxcblxuICAgICAgICBhcHA6IHtcbiAgICAgICAgICAgIGdldCBhY3RpdmVXaW5kb3coKTogc291cmNlZ3JhcGguV2luZG93IHwgdW5kZWZpbmVkIHtcbiAgICAgICAgICAgICAgICByZXR1cm4gd2luZG93cy5nZXRBY3RpdmUoKVxuICAgICAgICAgICAgfSxcbiAgICAgICAgICAgIGdldCB3aW5kb3dzKCk6IHNvdXJjZWdyYXBoLldpbmRvd1tdIHtcbiAgICAgICAgICAgICAgICByZXR1cm4gd2luZG93cy5nZXRBbGwoKVxuICAgICAgICAgICAgfSxcbiAgICAgICAgICAgIGNyZWF0ZVBhbmVsVmlldzogaWQgPT4gdmlld3MuY3JlYXRlUGFuZWxWaWV3KGlkKSxcbiAgICAgICAgfSxcblxuICAgICAgICB3b3Jrc3BhY2U6IHtcbiAgICAgICAgICAgIGdldCB0ZXh0RG9jdW1lbnRzKCk6IHNvdXJjZWdyYXBoLlRleHREb2N1bWVudFtdIHtcbiAgICAgICAgICAgICAgICByZXR1cm4gZG9jdW1lbnRzLmdldEFsbCgpXG4gICAgICAgICAgICB9LFxuICAgICAgICAgICAgb25EaWRPcGVuVGV4dERvY3VtZW50OiBkb2N1bWVudHMub25EaWRPcGVuVGV4dERvY3VtZW50LFxuICAgICAgICB9LFxuXG4gICAgICAgIGNvbmZpZ3VyYXRpb246IHtcbiAgICAgICAgICAgIGdldDogKCkgPT4gY29uZmlndXJhdGlvbi5nZXQoKSxcbiAgICAgICAgICAgIHN1YnNjcmliZTogbmV4dCA9PiBjb25maWd1cmF0aW9uLnN1YnNjcmliZShuZXh0KSxcbiAgICAgICAgfSxcblxuICAgICAgICBsYW5ndWFnZXM6IHtcbiAgICAgICAgICAgIHJlZ2lzdGVySG92ZXJQcm92aWRlcjogKHNlbGVjdG9yLCBwcm92aWRlcikgPT4gbGFuZ3VhZ2VGZWF0dXJlcy5yZWdpc3RlckhvdmVyUHJvdmlkZXIoc2VsZWN0b3IsIHByb3ZpZGVyKSxcbiAgICAgICAgICAgIHJlZ2lzdGVyRGVmaW5pdGlvblByb3ZpZGVyOiAoc2VsZWN0b3IsIHByb3ZpZGVyKSA9PlxuICAgICAgICAgICAgICAgIGxhbmd1YWdlRmVhdHVyZXMucmVnaXN0ZXJEZWZpbml0aW9uUHJvdmlkZXIoc2VsZWN0b3IsIHByb3ZpZGVyKSxcbiAgICAgICAgICAgIHJlZ2lzdGVyVHlwZURlZmluaXRpb25Qcm92aWRlcjogKHNlbGVjdG9yLCBwcm92aWRlcikgPT5cbiAgICAgICAgICAgICAgICBsYW5ndWFnZUZlYXR1cmVzLnJlZ2lzdGVyVHlwZURlZmluaXRpb25Qcm92aWRlcihzZWxlY3RvciwgcHJvdmlkZXIpLFxuICAgICAgICAgICAgcmVnaXN0ZXJJbXBsZW1lbnRhdGlvblByb3ZpZGVyOiAoc2VsZWN0b3IsIHByb3ZpZGVyKSA9PlxuICAgICAgICAgICAgICAgIGxhbmd1YWdlRmVhdHVyZXMucmVnaXN0ZXJJbXBsZW1lbnRhdGlvblByb3ZpZGVyKHNlbGVjdG9yLCBwcm92aWRlciksXG4gICAgICAgICAgICByZWdpc3RlclJlZmVyZW5jZVByb3ZpZGVyOiAoc2VsZWN0b3IsIHByb3ZpZGVyKSA9PlxuICAgICAgICAgICAgICAgIGxhbmd1YWdlRmVhdHVyZXMucmVnaXN0ZXJSZWZlcmVuY2VQcm92aWRlcihzZWxlY3RvciwgcHJvdmlkZXIpLFxuICAgICAgICB9LFxuXG4gICAgICAgIHNlYXJjaDoge1xuICAgICAgICAgICAgcmVnaXN0ZXJRdWVyeVRyYW5zZm9ybWVyOiBwcm92aWRlciA9PiBzZWFyY2gucmVnaXN0ZXJRdWVyeVRyYW5zZm9ybWVyKHByb3ZpZGVyKSxcbiAgICAgICAgfSxcblxuICAgICAgICBjb21tYW5kczoge1xuICAgICAgICAgICAgcmVnaXN0ZXJDb21tYW5kOiAoY29tbWFuZCwgY2FsbGJhY2spID0+IGNvbW1hbmRzLnJlZ2lzdGVyQ29tbWFuZCh7IGNvbW1hbmQsIGNhbGxiYWNrIH0pLFxuICAgICAgICAgICAgZXhlY3V0ZUNvbW1hbmQ6IChjb21tYW5kLCAuLi5hcmdzKSA9PiBjb21tYW5kcy5leGVjdXRlQ29tbWFuZChjb21tYW5kLCBhcmdzKSxcbiAgICAgICAgfSxcblxuICAgICAgICBpbnRlcm5hbDoge1xuICAgICAgICAgICAgc3luYyxcbiAgICAgICAgICAgIHVwZGF0ZUNvbnRleHQ6IHVwZGF0ZXMgPT4gY29udGV4dC51cGRhdGVDb250ZXh0KHVwZGF0ZXMpLFxuICAgICAgICAgICAgc291cmNlZ3JhcGhVUkw6IG5ldyBVUkkoaW5pdERhdGEuc291cmNlZ3JhcGhVUkwpLFxuICAgICAgICAgICAgY2xpZW50QXBwbGljYXRpb246IGluaXREYXRhLmNsaWVudEFwcGxpY2F0aW9uLFxuICAgICAgICB9LFxuICAgIH1cbn1cbiJdfQ==