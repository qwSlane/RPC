using Google.Protobuf.WellKnownTypes;
using Transport;
using Types;

namespace Client.RecordService;  

public interface IRecordClient
{ 
   Task SetNewRecord(Record record);
   Task<Level> GetBestN(BestLevelCount n);
}

public class RecordClient : IRecordClient
{
   public async Task SetNewRecord(Record args)
   {
      Request request = new Request
      {
         Args = Any.Pack(args),
         Method = 0,
      };
      
      var result = await RpcClient.Instance.Invoke(request);

      if (String.IsNullOrEmpty(result.Error) == false)
      {
         throw new Exception(result.Error);
      }
   }

   public async Task<Level> GetBestN(BestLevelCount args)
   {
      Request request = new Request
      {
         Args = Any.Pack(args),
         Method = 1,
      };
      
     var result = await RpcClient.Instance.Invoke(request);
     if (String.IsNullOrEmpty(result.Error) == false)
     {
        throw new Exception(result.Error);
     }
     return result.Result.Unpack<Level>();
   }
}