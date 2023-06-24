using System.Net.WebSockets;
using Google.Protobuf;
using Transport;

namespace Client;

public class RpcClient
{
    private ClientWebSocket _ws;

    private static RpcClient _RpcClient;

    public static RpcClient Instance
    {
        get
        {
            if (_RpcClient == null)
            {
                _RpcClient = new RpcClient();
            }
            return _RpcClient;
        }
    }

    private RpcClient()
    {
        _ws = new ClientWebSocket();
    }

    public async Task Connect(string uri, string token)
    {
        _ws.Options.SetRequestHeader("Authorization", token);
        await _ws.ConnectAsync(new Uri(uri), CancellationToken.None);
    } 


    public async Task<Response> Invoke(Request request)
    {
        var data = request.ToByteArray();
        await _ws.SendAsync(data, WebSocketMessageType.Binary,true, CancellationToken.None);

        var buffer = new byte[1024];
        var result = await _ws.ReceiveAsync(new ArraySegment<byte>(buffer), CancellationToken.None);
        
        var messageBytes = new byte[result.Count];
        Array.Copy(buffer, messageBytes, result.Count);
        return Response.Parser.ParseFrom(messageBytes);
    }

    public async void CloseConnection()
    {
        await _ws.CloseAsync(WebSocketCloseStatus.NormalClosure, "Connection closed" ,CancellationToken.None);
    }
}