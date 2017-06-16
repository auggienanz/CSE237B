% CSE 237B, Final Project
% August Nanz

data = csvread('./latency.csv');
time = (data(:,1) - data(1,1))./1e9;
latency = data(:,2)./1e6;
y_s = data(:,3)./1e6;
y_up = data(:,4)./1e6;

figure;
plot(time, latency);
hold on
plot(time, y_s);
plot(time,y_up);
legend('Raw','Smoothed','Upper Bound');
title('Smoothed TCP Latency');
xlabel('Time (seconds)');
ylabel('Latency (ms)');
xlim([0 300]);
ylim([30 130]);

%% For testing
alpha = 0.9;
beta = 0.9;
kappa = 0.1;

y_s = zeros(size(time,1),1);
y_s(1) = latency(1);
y_up = zeros(size(time,1),1);
y_up(1) = latency(1);
y_var = 0;

for i = 2:size(time,1)
    y_var = (1 - beta) * y_var + beta * abs(y_s(i-1) - latency(i));
    y_s(i) = (1 - alpha) * y_s(i-1) + alpha * latency(i);
    y_up(i) = y_s(i) + kappa*y_var;
end
figure;
gca.ColorOrderIndex = 1;
plot(time, latency);
hold on
plot(time, y_s);
plot(time,y_up);
legend('Raw','Smoothed','Upper Bound');
title('Alternate Parameters Smoothed TCP Latency');
xlabel('Time (seconds)');
ylabel('Latency (ms)');
xlim([0 300]);
ylim([30 130]);