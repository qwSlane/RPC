
// Generated code.
// DO NOT EDIT.

using Client;
using Google.Protobuf.WellKnownTypes;
using Transport;
using Types;

public class RecordsClient
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
         Method = 0,
      };
      
      var result = await RpcClient.Instance.Invoke(request);

      if (String.IsNullOrEmpty(result.Error) == false)
      {
         throw new Exception(result.Error);
      }

      return result.Result.Unpack<Level>();
   }
}

	