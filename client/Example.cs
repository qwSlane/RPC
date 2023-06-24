using System.Net.WebSockets;
using Client;
using Client.RecordService;
using Types;

class Program
{
    static async Task Main()
    {
        using WebSocket ws = new ClientWebSocket();
        Record record = new Record
        {
            Level = 1,
            Username = "CsUser",
            Score = 12000,
        };

        await RpcClient.Instance.Connect("ws://localhost:8080/ws",
            "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2ODc2MTA5MzgsInVzZXJuYW1lIjoiUm9iZXJ0IFBvbHNvbiJ9.jsLpTDfejPgZVOIanjbP1Xqy7xAzMX1I31NXZt9Zm80");

        BestLevelCount best = new BestLevelCount
        {
            Count = 3,
            Level = 1
        };
        
        RecordClient client = new();

        await client.GetBestN(best);
        var result = await client.GetBestN(best);

        foreach (var score in result.Scores)
        {
            Console.WriteLine($"{score.Username} : {score.Score}");
        }
        
        RpcClient.Instance.CloseConnection();

    }
}