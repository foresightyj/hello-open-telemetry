using OpenTelemetry;
using OpenTelemetry.Instrumentation.AspNetCore;
using OpenTelemetry.Trace;

namespace dotnet_core_gateway
{
    public class Program
    {
        public static void Main(string[] args)
        {
            var builder = WebApplication.CreateBuilder(args);

            builder.Services.AddControllers();
            // Learn more about configuring Swagger/OpenAPI at https://aka.ms/aspnetcore/swashbuckle
            builder.Services.AddEndpointsApiExplorer();
            builder.Services.AddSwaggerGen();
            builder.Services.AddOpenTelemetryTracing((builder) =>
               {
                   builder
                   .AddSource("weather-gateway")
                  .AddAspNetCoreInstrumentation()
                  .AddHttpClientInstrumentation((options) =>
                    options.Filter = (httpRequestMessage) =>
                    {
                        return httpRequestMessage.Method.Equals(HttpMethod.Get);
                    })
                  .AddJaegerExporter()
                  .AddConsoleExporter();
               });

            builder.Services.AddHttpClient();

            builder.Services.Configure<AspNetCoreInstrumentationOptions>(options =>
            {
                options.Filter = (httpContext) =>
                {
                    // only collect telemetry about HTTP GET requests
                    return httpContext.Request.Method.Equals("GET");
                };
            });

            var app = builder.Build();

            // Configure the HTTP request pipeline.
            if (app.Environment.IsDevelopment())
            {
                app.UseSwagger();
                app.UseSwaggerUI();
            }

            app.UseAuthorization();

            app.MapControllers();

            app.Run();
        }
    }
}